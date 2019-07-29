package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"reflect"
	"regexp"

	"github.com/golang/glog"

	admission_v1beta1 "k8s.io/api/admission/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	shipper "github.com/bookingcom/shipper/pkg/apis/shipper/v1alpha1"
	clientset "github.com/bookingcom/shipper/pkg/client/clientset/versioned"
	shippererrors "github.com/bookingcom/shipper/pkg/errors"
	"github.com/bookingcom/shipper/pkg/util/rolloutblock"
	kubeclient "k8s.io/api/admission/v1beta1"
)

type Webhook struct {
	shipperClientset clientset.Interface
	bindAddr string
	bindPort string

	tlsCertFile       string
	tlsPrivateKeyFile string
}

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

func NewWebhook(bindAddr, bindPort, tlsPrivateKeyFile, tlsCertFile string, shipperClientset clientset.Interface) *Webhook {
	return &Webhook{
		shipperClientset:  shipperClientset,
		bindAddr:          bindAddr,
		bindPort:          bindPort,
		tlsPrivateKeyFile: tlsPrivateKeyFile,
		tlsCertFile:       tlsCertFile,
	}
}

func (c *Webhook) Run(stopCh <-chan struct{}) {
	addr := c.bindAddr + ":" + c.bindPort
	mux := c.initializeHandlers()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		var serverError error
		if c.tlsCertFile == "" || c.tlsPrivateKeyFile == "" {
			serverError = server.ListenAndServe()
		} else {
			serverError = server.ListenAndServeTLS(c.tlsCertFile, c.tlsPrivateKeyFile)
		}

		if serverError != nil && serverError != http.ErrServerClosed {
			glog.Fatalf("failed to start shipper-webhook: %v", serverError)
		}
	}()

	glog.V(2).Info("Started the WebHook")

	<-stopCh

	glog.V(2).Info("Shutting down the WebHook")

	if err := server.Shutdown(context.Background()); err != nil {
		glog.Errorf(`HTTP server Shutdown: %v`, err)
	}
}

func (c *Webhook) initializeHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", adaptHandler(c.validateHandlerFunc))
	return mux
}

// adaptHandler wraps an admission review function to be consumed through HTTP.
func adaptHandler(handler func(*admission_v1beta1.AdmissionReview) *admission_v1beta1.AdmissionResponse) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		if r.Body != nil {
			if data, err := ioutil.ReadAll(r.Body); err == nil {
				body = data
			}
		}

		if len(body) == 0 {
			http.Error(w, "empty body", http.StatusBadRequest)
			return
		}

		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, "Invalid content-type", http.StatusUnsupportedMediaType)
			return
		}

		if mediaType != "application/json" {
			http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
			return
		}

		var admissionResponse *admission_v1beta1.AdmissionResponse
		ar := admission_v1beta1.AdmissionReview{}
		if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
			admissionResponse = &admission_v1beta1.AdmissionResponse{
				Result: &meta_v1.Status{
					Message: err.Error(),
				},
			}
		} else {
			admissionResponse = handler(&ar)
		}

		admissionReview := admission_v1beta1.AdmissionReview{}
		if admissionResponse != nil {
			admissionReview.Response = admissionResponse
			if ar.Request != nil {
				admissionReview.Response.UID = ar.Request.UID
			}
		}

		resp, err := json.Marshal(admissionReview)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(resp); err != nil {
			http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func (c *Webhook) validateHandlerFunc(review *admission_v1beta1.AdmissionReview) *admission_v1beta1.AdmissionResponse {
	request := review.Request
	var err error

	switch request.Kind.Kind {
	case "Application":
		var application shipper.Application
		err = json.Unmarshal(request.Object.Raw, &application)
		if err == nil {
			err = c.validateApplication(request, application)
		}
	case "Release":
		var release shipper.Release
		err = json.Unmarshal(request.Object.Raw, &release)
		if err == nil {
			err = c.validateRelease(request, release)
		}
	case "Cluster":
		var cluster shipper.Cluster
		err = json.Unmarshal(request.Object.Raw, &cluster)
	case "InstallationTarget":
		var installationTarget shipper.InstallationTarget
		err = json.Unmarshal(request.Object.Raw, &installationTarget)
	case "CapacityTarget":
		var capacityTarget shipper.CapacityTarget
		err = json.Unmarshal(request.Object.Raw, &capacityTarget)
	case "TrafficTarget":
		var trafficTarget shipper.TrafficTarget
		err = json.Unmarshal(request.Object.Raw, &trafficTarget)
	case "RolloutBlock":
		var rolloutBlock shipper.RolloutBlock
		err = json.Unmarshal(request.Object.Raw, &rolloutBlock)
	}

	if err != nil {
		return &admission_v1beta1.AdmissionResponse{
			Result: &meta_v1.Status{
				Message: err.Error(),
			},
		}
	}

	return &admission_v1beta1.AdmissionResponse{
		Allowed: true,
	}
}

func (c *Webhook) validateRelease(request *admission_v1beta1.AdmissionRequest, release shipper.Release) error {
	var err error
	overrideRBs := rolloutblock.NewOverride(release.Annotations[shipper.RolloutBlocksOverrideAnnotation])
	err = c.validateOverrideRolloutBlockAnnotation(overrideRBs, release.Namespace)
	if err != nil {
		return err
	}

	if request.Operation == kubeclient.Create {
		err = c.processRolloutBlocks(release.Namespace, overrideRBs)
	}

	if request.Operation == kubeclient.Update {
		var oldRelease shipper.Release
		err = json.Unmarshal(request.OldObject.Raw, &oldRelease)
		if err != nil {
			return nil
		}
		newOverrides := rolloutblock.NewOverride(release.Annotations[shipper.RolloutBlocksOverrideAnnotation])
		oldOverrides := rolloutblock.NewOverride(oldRelease.Annotations[shipper.RolloutBlocksOverrideAnnotation])
		subtractedOverrides := oldOverrides.Diff(newOverrides)
		if len(subtractedOverrides) == 0 {
			// no overrides were removed from annotation. validating release:
			err = c.processRolloutBlocks(release.Namespace, overrideRBs)
		} else {
			// rolloutblock override were changed, making sure spec was not changed
			if !reflect.DeepEqual(release.Spec, oldRelease.Spec) {
				err = fmt.Errorf("releases Spec was changed %v", release)
			}
		}
	}

	return err
}

func (c *Webhook) validateApplication(request *admission_v1beta1.AdmissionRequest, application shipper.Application) error {
	var err error
	overrideRBs := rolloutblock.NewOverride(application.Annotations[shipper.RolloutBlocksOverrideAnnotation])
	err = c.validateOverrideRolloutBlockAnnotation(overrideRBs, application.Namespace)
	if err != nil {
		return err
	}

	if request.Operation == kubeclient.Create {
		err = c.processRolloutBlocks(application.Namespace, overrideRBs)
	}

	return err
}

func (c *Webhook) processRolloutBlocks(namespace string, overrideRBs rolloutblock.Override) error {
	existingRBs := rolloutblock.NewOverrideFromRolloutBlocks(c.existingRolloutBlocks(namespace))

	nonOverriddenRBs := existingRBs.Diff(overrideRBs)
	if len(nonOverriddenRBs) > 0 {
		return shippererrors.NewRolloutBlockError(nonOverriddenRBs.String())
	}

	return nil
}

func (c *Webhook) existingRolloutBlocks(namespace string) []*shipper.RolloutBlock {
	var nsRBs, gbRBs []*shipper.RolloutBlock
	if nsRBList, err := c.shipperClientset.ShipperV1alpha1().RolloutBlocks(namespace).List(meta_v1.ListOptions{}); err == nil {
		for _, item := range nsRBList.Items {
			nsRBs = append(nsRBs, &item)
		}
	}
	if gbRBList, err := c.shipperClientset.ShipperV1alpha1().RolloutBlocks(shipper.GlobalRolloutBlockNamespace).List(meta_v1.ListOptions{}); err == nil {
		for _, item := range gbRBList.Items {
			gbRBs = append(gbRBs, &item)
		}
	}
	return append(nsRBs, gbRBs...)
}

func (c *Webhook) validateOverrideRolloutBlockAnnotation(overrideRbs rolloutblock.Override, namespace string) error {
	if len(overrideRbs) == 0 {
		return nil
	}

	re := regexp.MustCompile("^[a-zA-Z0-9/-]+/[a-zA-Z0-9/-]+$")

	for item := range overrideRbs {
		if !re.MatchString(item) {
			return shippererrors.NewInvalidRolloutBlockOverrideError(item)
		}
	}

	existingRbs := rolloutblock.NewOverrideFromRolloutBlocks(c.existingRolloutBlocks(namespace))
	nonExistingRbs := overrideRbs.Diff(existingRbs)

	if len(nonExistingRbs) > 0 {
		return shippererrors.NewInvalidRolloutBlockOverrideError(nonExistingRbs.String())
	}

	return nil
}
