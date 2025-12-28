package schemas

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Response struct {
	Status bool `json:"status"`
	Body any `json:"body"`
	Message string `json:"message"`
}

type ErrorStruc struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *ErrorStruc) Error() string {
	return e.Message
}

func ErrorResponse(code int, message string, err error) *ErrorStruc {
	return &ErrorStruc{
		StatusCode: code,
		Message:    message,
		Err:        err,
	}
}

func HandleError(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	if errResp, ok := err.(*ErrorStruc); ok {
		log.Err(err).Msgf("Error: %s", errResp.Err.Error())
		return ctx.Status(errResp.StatusCode).JSON(Response{
			Status:  false,
			Body:    nil,
			Message: errResp.Message,
		})
	}

	log.Err(err).Msgf("Error: %s", err.Error())
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Status:  false,
		Body:    nil,
		Message: "Error interno",
	})
}

func HandlerErrorGrpc(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		log.Error().Err(err).Msg("Non-gRPC Error")
		return ErrorResponse(fiber.StatusInternalServerError, "Error de sistema", err)
	}

	var httpStatus int
	switch st.Code() {
	case codes.InvalidArgument:
		httpStatus = fiber.StatusBadRequest
	case codes.Unauthenticated:
		httpStatus = fiber.StatusUnauthorized
	case codes.PermissionDenied:
		httpStatus = fiber.StatusForbidden
	case codes.NotFound:
		httpStatus = fiber.StatusNotFound
	case codes.AlreadyExists:
		httpStatus = fiber.StatusConflict
	case codes.DeadlineExceeded:
		httpStatus = fiber.StatusGatewayTimeout
	case codes.Unavailable:
		httpStatus = fiber.StatusServiceUnavailable
	default:
		httpStatus = fiber.StatusInternalServerError
	}

	log.Error().
		Int("grpc_code", int(st.Code())).
		Int("http_code", httpStatus).
		Msgf("gRPC Error: %s", st.Message())

	return ErrorResponse(httpStatus, st.Message(), err)
}