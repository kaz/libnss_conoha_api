package main

// #include <nss.h>
// #include <netdb.h>
// #include <stdlib.h>
// #define PTR_SIZE sizeof(char*)
import "C"
import "unsafe"

//export _nss_conoha_gethostbyname2_r
func _nss_conoha_gethostbyname2_r(name *C.char, af C.int, ret *C.struct_hostent, buf *C.char, buflen C.size_t, errnop *C.int, h_errnop *C.int) C.enum_nss_status {
	*errnop = 0
	*h_errnop = 0

	if af == C.AF_INET {
		ret.h_length = 4
	} else if af == C.AF_INET6 {
		ret.h_length = 16
	} else {
		return C.NSS_STATUS_NOTFOUND
	}

	result := LookupInstance(C.GoString(name), af == C.AF_INET6)
	if result == nil {
		return C.NSS_STATUS_NOTFOUND
	}

	ret.h_name = name
	ret.h_aliases = (**C.char)(C.calloc(C.PTR_SIZE, 1))
	ret.h_addrtype = af

	ptr := C.calloc(C.PTR_SIZE, C.ulong(len(result)+1))
	ret.h_addr_list = (**C.char)(ptr)

	for i, addr := range result {
		*(**C.char)(unsafe.Pointer(uintptr(ptr) + uintptr(i*C.PTR_SIZE))) = C.CString(addr)
	}

	return C.NSS_STATUS_SUCCESS
}

//export _nss_conoha_gethostbyname_r
func _nss_conoha_gethostbyname_r(name *C.char, ret *C.struct_hostent, buf *C.char, buflen C.size_t, errnop *C.int, h_errnop *C.int) C.enum_nss_status {
	return _nss_conoha_gethostbyname2_r(name, C.AF_INET, ret, buf, buflen, errnop, h_errnop)
}
