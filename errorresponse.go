package hg

// ErrorResponse returns the appropriate response type for the given status-code & meta.
//
// And note that any Mercury Protocol response type other that “20 SUCCESS” is considered an error.
//
// Example With Bad Request
//
// So, for example, this:
//
//	hg.ErrorResponse(59, "")
//
// Would return:
//
//	hg.ResponseBadRequest{meta:""}
//
// Example With Temporary Failure
//
// And, for example, this:
//
//	hg.ErrorResponse(40, "we seem to be experiencing some technical difficulties")
//
// Would return:
//
//	hg.ResponseTemporaryFailure{meta:"we seem to be experiencing some technical difficulties"}
//
// Example With Success
//
// Althugh note that calling with the a status-code of 20 (i.e., the status code for Success) would return nil.
// So, for example, this:
//
//	hg.ErrorResponse(20, "text/gemini")
//
// Would return
//
//	nil
//
// Type Switch
//
// This is useful with type switches. For example:
//
//	func callMercury(rr hg.ResponseReader, r hg.Request) {
//		
//		// ...
//		
//		p, err := ioutil.ReadAll(rr)
//		
//		if nil != err {
//			switch casted: err.(type) {
//			case hg.ResponseInput:
//				//@TODO
//			case hg.ResponseSensitiveInput:
//				//@TODO
//
//			case hg.ResponseRedirectTemporary:
//				//@TODO
//			case hg.ResponseRedirectPermanent:
//				//@TODO
//
//			case hg.ResponseTemporaryFailure:
//				//@TODO
//			case hg.ResponseServerUnavailable:
//				//@TODO
//			case hg.ResponseCGIError:
//				//@TODO
//			case hg.ResponseProxyError:
//				//@TODO
//			case hg.ResponseSlowDown:
//				//@TODO
//
//			case hg.ResponsePermanentFailure:
//				//@TODO
//			case hg.ResponseNotFound :
//				//@TODO
//			case hg.ResponseGone:
//				//@TODO
//			case hg.ResponseProxyRequestRefused:
//				//@TODO
//			case hg.ResponseBadRequest:
//				//@TODO
//
//			case hg.UnknownResponse:
//				//@TODO
//
//			default:
//				//@TODO
//		}
//		
//		// ...
//		
//	}
func ErrorResponse(statuscode int, meta string) error {
	switch statuscode {
	case StatusInput:
		return ResponseInput{meta}
	case StatusSensitiveInput:
		return ResponseSensitiveInput{meta}

	case StatusSuccess:
		return nil

	case StatusRedirectTemporary:
		return ResponseRedirectTemporary{meta}
	case StatusRedirectPermanent:
		return ResponseRedirectPermanent{meta}

	case StatusTemporaryFailure:
		return ResponseTemporaryFailure{meta}
	case StatusServerUnavailable:
		return ResponseServerUnavailable{meta}
	case StatusCGIError:
		return ResponseCGIError{meta}
	case StatusProxyError:
		return ResponseProxyError{meta}
	case StatusSlowDown:
		return ResponseSlowDown{meta}

	case StatusPermanentFailure:
		return ResponsePermanentFailure{meta}
	case StatusNotFound:
		return ResponseNotFound{meta}
	case StatusGone:
		return ResponseGone{meta}
	case StatusProxyRequestRefused:
		return ResponseProxyRequestRefused{meta}
	case StatusBadRequest:
		return ResponseBadRequest{meta}

	default:
		return UnknownResponse{meta:meta, statusCode:statuscode}
	}
}
