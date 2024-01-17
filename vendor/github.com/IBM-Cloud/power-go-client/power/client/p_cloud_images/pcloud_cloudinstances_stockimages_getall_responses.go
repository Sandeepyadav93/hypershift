// Code generated by go-swagger; DO NOT EDIT.

package p_cloud_images

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/IBM-Cloud/power-go-client/power/models"
)

// PcloudCloudinstancesStockimagesGetallReader is a Reader for the PcloudCloudinstancesStockimagesGetall structure.
type PcloudCloudinstancesStockimagesGetallReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PcloudCloudinstancesStockimagesGetallReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPcloudCloudinstancesStockimagesGetallOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPcloudCloudinstancesStockimagesGetallBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewPcloudCloudinstancesStockimagesGetallUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewPcloudCloudinstancesStockimagesGetallForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewPcloudCloudinstancesStockimagesGetallNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewPcloudCloudinstancesStockimagesGetallInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images] pcloud.cloudinstances.stockimages.getall", response, response.Code())
	}
}

// NewPcloudCloudinstancesStockimagesGetallOK creates a PcloudCloudinstancesStockimagesGetallOK with default headers values
func NewPcloudCloudinstancesStockimagesGetallOK() *PcloudCloudinstancesStockimagesGetallOK {
	return &PcloudCloudinstancesStockimagesGetallOK{}
}

/*
PcloudCloudinstancesStockimagesGetallOK describes a response with status code 200, with default header values.

OK
*/
type PcloudCloudinstancesStockimagesGetallOK struct {
	Payload *models.Images
}

// IsSuccess returns true when this pcloud cloudinstances stockimages getall o k response has a 2xx status code
func (o *PcloudCloudinstancesStockimagesGetallOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this pcloud cloudinstances stockimages getall o k response has a 3xx status code
func (o *PcloudCloudinstancesStockimagesGetallOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this pcloud cloudinstances stockimages getall o k response has a 4xx status code
func (o *PcloudCloudinstancesStockimagesGetallOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this pcloud cloudinstances stockimages getall o k response has a 5xx status code
func (o *PcloudCloudinstancesStockimagesGetallOK) IsServerError() bool {
	return false
}

// IsCode returns true when this pcloud cloudinstances stockimages getall o k response a status code equal to that given
func (o *PcloudCloudinstancesStockimagesGetallOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the pcloud cloudinstances stockimages getall o k response
func (o *PcloudCloudinstancesStockimagesGetallOK) Code() int {
	return 200
}

func (o *PcloudCloudinstancesStockimagesGetallOK) Error() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallOK  %+v", 200, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallOK) String() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallOK  %+v", 200, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallOK) GetPayload() *models.Images {
	return o.Payload
}

func (o *PcloudCloudinstancesStockimagesGetallOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Images)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPcloudCloudinstancesStockimagesGetallBadRequest creates a PcloudCloudinstancesStockimagesGetallBadRequest with default headers values
func NewPcloudCloudinstancesStockimagesGetallBadRequest() *PcloudCloudinstancesStockimagesGetallBadRequest {
	return &PcloudCloudinstancesStockimagesGetallBadRequest{}
}

/*
PcloudCloudinstancesStockimagesGetallBadRequest describes a response with status code 400, with default header values.

Bad Request
*/
type PcloudCloudinstancesStockimagesGetallBadRequest struct {
	Payload *models.Error
}

// IsSuccess returns true when this pcloud cloudinstances stockimages getall bad request response has a 2xx status code
func (o *PcloudCloudinstancesStockimagesGetallBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this pcloud cloudinstances stockimages getall bad request response has a 3xx status code
func (o *PcloudCloudinstancesStockimagesGetallBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this pcloud cloudinstances stockimages getall bad request response has a 4xx status code
func (o *PcloudCloudinstancesStockimagesGetallBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this pcloud cloudinstances stockimages getall bad request response has a 5xx status code
func (o *PcloudCloudinstancesStockimagesGetallBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this pcloud cloudinstances stockimages getall bad request response a status code equal to that given
func (o *PcloudCloudinstancesStockimagesGetallBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the pcloud cloudinstances stockimages getall bad request response
func (o *PcloudCloudinstancesStockimagesGetallBadRequest) Code() int {
	return 400
}

func (o *PcloudCloudinstancesStockimagesGetallBadRequest) Error() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallBadRequest  %+v", 400, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallBadRequest) String() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallBadRequest  %+v", 400, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallBadRequest) GetPayload() *models.Error {
	return o.Payload
}

func (o *PcloudCloudinstancesStockimagesGetallBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPcloudCloudinstancesStockimagesGetallUnauthorized creates a PcloudCloudinstancesStockimagesGetallUnauthorized with default headers values
func NewPcloudCloudinstancesStockimagesGetallUnauthorized() *PcloudCloudinstancesStockimagesGetallUnauthorized {
	return &PcloudCloudinstancesStockimagesGetallUnauthorized{}
}

/*
PcloudCloudinstancesStockimagesGetallUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type PcloudCloudinstancesStockimagesGetallUnauthorized struct {
	Payload *models.Error
}

// IsSuccess returns true when this pcloud cloudinstances stockimages getall unauthorized response has a 2xx status code
func (o *PcloudCloudinstancesStockimagesGetallUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this pcloud cloudinstances stockimages getall unauthorized response has a 3xx status code
func (o *PcloudCloudinstancesStockimagesGetallUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this pcloud cloudinstances stockimages getall unauthorized response has a 4xx status code
func (o *PcloudCloudinstancesStockimagesGetallUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this pcloud cloudinstances stockimages getall unauthorized response has a 5xx status code
func (o *PcloudCloudinstancesStockimagesGetallUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this pcloud cloudinstances stockimages getall unauthorized response a status code equal to that given
func (o *PcloudCloudinstancesStockimagesGetallUnauthorized) IsCode(code int) bool {
	return code == 401
}

// Code gets the status code for the pcloud cloudinstances stockimages getall unauthorized response
func (o *PcloudCloudinstancesStockimagesGetallUnauthorized) Code() int {
	return 401
}

func (o *PcloudCloudinstancesStockimagesGetallUnauthorized) Error() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallUnauthorized  %+v", 401, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallUnauthorized) String() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallUnauthorized  %+v", 401, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallUnauthorized) GetPayload() *models.Error {
	return o.Payload
}

func (o *PcloudCloudinstancesStockimagesGetallUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPcloudCloudinstancesStockimagesGetallForbidden creates a PcloudCloudinstancesStockimagesGetallForbidden with default headers values
func NewPcloudCloudinstancesStockimagesGetallForbidden() *PcloudCloudinstancesStockimagesGetallForbidden {
	return &PcloudCloudinstancesStockimagesGetallForbidden{}
}

/*
PcloudCloudinstancesStockimagesGetallForbidden describes a response with status code 403, with default header values.

Forbidden
*/
type PcloudCloudinstancesStockimagesGetallForbidden struct {
	Payload *models.Error
}

// IsSuccess returns true when this pcloud cloudinstances stockimages getall forbidden response has a 2xx status code
func (o *PcloudCloudinstancesStockimagesGetallForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this pcloud cloudinstances stockimages getall forbidden response has a 3xx status code
func (o *PcloudCloudinstancesStockimagesGetallForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this pcloud cloudinstances stockimages getall forbidden response has a 4xx status code
func (o *PcloudCloudinstancesStockimagesGetallForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this pcloud cloudinstances stockimages getall forbidden response has a 5xx status code
func (o *PcloudCloudinstancesStockimagesGetallForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this pcloud cloudinstances stockimages getall forbidden response a status code equal to that given
func (o *PcloudCloudinstancesStockimagesGetallForbidden) IsCode(code int) bool {
	return code == 403
}

// Code gets the status code for the pcloud cloudinstances stockimages getall forbidden response
func (o *PcloudCloudinstancesStockimagesGetallForbidden) Code() int {
	return 403
}

func (o *PcloudCloudinstancesStockimagesGetallForbidden) Error() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallForbidden  %+v", 403, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallForbidden) String() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallForbidden  %+v", 403, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallForbidden) GetPayload() *models.Error {
	return o.Payload
}

func (o *PcloudCloudinstancesStockimagesGetallForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPcloudCloudinstancesStockimagesGetallNotFound creates a PcloudCloudinstancesStockimagesGetallNotFound with default headers values
func NewPcloudCloudinstancesStockimagesGetallNotFound() *PcloudCloudinstancesStockimagesGetallNotFound {
	return &PcloudCloudinstancesStockimagesGetallNotFound{}
}

/*
PcloudCloudinstancesStockimagesGetallNotFound describes a response with status code 404, with default header values.

Not Found
*/
type PcloudCloudinstancesStockimagesGetallNotFound struct {
	Payload *models.Error
}

// IsSuccess returns true when this pcloud cloudinstances stockimages getall not found response has a 2xx status code
func (o *PcloudCloudinstancesStockimagesGetallNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this pcloud cloudinstances stockimages getall not found response has a 3xx status code
func (o *PcloudCloudinstancesStockimagesGetallNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this pcloud cloudinstances stockimages getall not found response has a 4xx status code
func (o *PcloudCloudinstancesStockimagesGetallNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this pcloud cloudinstances stockimages getall not found response has a 5xx status code
func (o *PcloudCloudinstancesStockimagesGetallNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this pcloud cloudinstances stockimages getall not found response a status code equal to that given
func (o *PcloudCloudinstancesStockimagesGetallNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the pcloud cloudinstances stockimages getall not found response
func (o *PcloudCloudinstancesStockimagesGetallNotFound) Code() int {
	return 404
}

func (o *PcloudCloudinstancesStockimagesGetallNotFound) Error() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallNotFound  %+v", 404, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallNotFound) String() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallNotFound  %+v", 404, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *PcloudCloudinstancesStockimagesGetallNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPcloudCloudinstancesStockimagesGetallInternalServerError creates a PcloudCloudinstancesStockimagesGetallInternalServerError with default headers values
func NewPcloudCloudinstancesStockimagesGetallInternalServerError() *PcloudCloudinstancesStockimagesGetallInternalServerError {
	return &PcloudCloudinstancesStockimagesGetallInternalServerError{}
}

/*
PcloudCloudinstancesStockimagesGetallInternalServerError describes a response with status code 500, with default header values.

Internal Server Error
*/
type PcloudCloudinstancesStockimagesGetallInternalServerError struct {
	Payload *models.Error
}

// IsSuccess returns true when this pcloud cloudinstances stockimages getall internal server error response has a 2xx status code
func (o *PcloudCloudinstancesStockimagesGetallInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this pcloud cloudinstances stockimages getall internal server error response has a 3xx status code
func (o *PcloudCloudinstancesStockimagesGetallInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this pcloud cloudinstances stockimages getall internal server error response has a 4xx status code
func (o *PcloudCloudinstancesStockimagesGetallInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this pcloud cloudinstances stockimages getall internal server error response has a 5xx status code
func (o *PcloudCloudinstancesStockimagesGetallInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this pcloud cloudinstances stockimages getall internal server error response a status code equal to that given
func (o *PcloudCloudinstancesStockimagesGetallInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the pcloud cloudinstances stockimages getall internal server error response
func (o *PcloudCloudinstancesStockimagesGetallInternalServerError) Code() int {
	return 500
}

func (o *PcloudCloudinstancesStockimagesGetallInternalServerError) Error() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallInternalServerError  %+v", 500, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallInternalServerError) String() string {
	return fmt.Sprintf("[GET /pcloud/v1/cloud-instances/{cloud_instance_id}/stock-images][%d] pcloudCloudinstancesStockimagesGetallInternalServerError  %+v", 500, o.Payload)
}

func (o *PcloudCloudinstancesStockimagesGetallInternalServerError) GetPayload() *models.Error {
	return o.Payload
}

func (o *PcloudCloudinstancesStockimagesGetallInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
