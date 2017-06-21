package ListaModel

import "html/template"

//SListaVista estructura Auxiliar para enviar listas a la vista
type SListaVista struct {
	ID      string
	SEstado bool
	SMsj    string
	SGrupo  template.HTML
}
