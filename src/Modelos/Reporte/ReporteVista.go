package reporte

import "time"

// ReporteVista es un elemento que representa la estructura reporte en la vista
type ReporteVista struct {
	CodigoBarraTicket   CodigoBarraTicketVista
	CodigoBarraSurtidor CodigoBarraSurtidorVista
	TimeIn              TimeInVista
	TimeOut             TimeOutVista
	DuracionM           DuracionMVista
	Respuesta           RespuestaVista
	Estado              bool
	Error               string
	Mensaje             string
	TimerOn             bool
	Concluido           bool
}

// CodigoBarraTicketVista estructura que representa el miembro "CodigoBarraTicket" de la estructura reporte en la vista.
type CodigoBarraTicketVista struct {
	CodigoBarraTicket string //Valor
	Error             string
	Mensaje           string
	Estado            bool
}

// CodigoBarraSurtidorVista estructura que representa el miembro "CodigoBarraSurtidor" de la estructura reporte en la vista.
type CodigoBarraSurtidorVista struct {
	CodigoBarraSurtidor string //Valor
	Error               string
	Mensaje             string
	Estado              bool
}

// TimeInVista estructura que representa el miembro "TimeIn" de la estructura reporte en la vista.
type TimeInVista struct {
	TimeIn  time.Time //Valor
	Error   string
	Mensaje string
	Estado  bool
}

// TimeOutVista estructura que representa el miembro "TimeOutVista" de la estructura reporte en la vista.
type TimeOutVista struct {
	TimeOut time.Time //Valor
	Error   string
	Mensaje string
	Estado  bool
}

// DuracionMVista estructura que representa el miembro "DuracionM" de la estructura reporte en la vista.
type DuracionMVista struct {
	DuracionM int64 //Valor
	Error     string
	Mensaje   string
	Estado    bool
}

// RespuestaVista estructura que representa el miembro "Respuesta" de la estructura reporte en la vista.
type RespuestaVista struct {
	Respuesta string //Valor
	Error     string
	Mensaje   string
	Estado    bool
}
