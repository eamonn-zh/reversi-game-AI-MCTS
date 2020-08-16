from ctypes import *


class PiecePosition(Structure):
    _fields_ = [('X', c_int), ('Y', c_int)]


# go type
# from https://stackoverflow.com/questions/56586267/calling-go-from-python
class GoSlice(Structure):
    _fields_ = [("data", POINTER(c_void_p)), ("len", c_longlong), ("cap", c_longlong)]