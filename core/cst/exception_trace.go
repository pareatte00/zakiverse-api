package cst

type ExceptionTrace string

const (
	ExceptionTraceHandler    ExceptionTrace = "[Handler]"
	ExceptionTraceService    ExceptionTrace = "[Service]"
	ExceptionTraceRepository ExceptionTrace = "[Repository]"
	ExceptionTraceUtils      ExceptionTrace = "[Utils]"
)
