package domain

import (
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type UserController struct {
	service Service
	tracer  trace.Tracer
}

func NewUserController(service Service, tracer trace.Tracer) UserController {
	return UserController{
		service: service,
		tracer:  tracer,
	}
}

func (controller *UserController) Get(w http.ResponseWriter, req *http.Request) {
	_, span := controller.tracer.Start(req.Context(), "controller::Get")
	defer span.End(trace.WithStackTrace(true))

	userId := req.URL.Query().Get("id")
	user, _ := controller.service.Get(userId)
	setResponse(w, user, http.StatusOK)
}

func (controller *UserController) Add(w http.ResponseWriter, req *http.Request) {
	_, span := controller.tracer.Start(req.Context(), "controller::Add")
	defer span.End(trace.WithStackTrace(true))
	var user User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		span.SetStatus(otelcodes.Error, "request decode error")
		span.RecordError(err, trace.WithStackTrace(true))
		setResponse(w, err, http.StatusBadRequest)
		return
	}

	error := controller.service.Add(user)
	if error != nil {
		span.RecordError(error, trace.WithStackTrace(true))
		setResponse(w, error.Error(), http.StatusBadRequest)
	} else {
		span.AddEvent("user", trace.WithAttributes(attribute.String("Id", user.Id)))
		span.SetStatus(otelcodes.Ok, "user added")
		span.SetAttributes(attribute.String("Test", "Test"))
		setResponse(w, user, http.StatusCreated)
	}

}

func (controller *UserController) All(w http.ResponseWriter, req *http.Request) {
	_, span := controller.tracer.Start(req.Context(), "controller::All")
	defer span.End(trace.WithStackTrace(true))
	users, _ := controller.service.List()
	setResponse(w, users, http.StatusOK)
}

func setResponse(w http.ResponseWriter, output interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(output)
}
