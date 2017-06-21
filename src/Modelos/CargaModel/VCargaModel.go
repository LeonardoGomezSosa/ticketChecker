package CargaModel

//SCargaVista estructura Auxiliar para enviar listas a la vista
type SCargaVista struct {
	ID      string
	SEstado bool
	SMsj    string
	Error   []Errores
	//SGrupo  template.HTML
}

//Errores interfaz para capturar errores
type Errores struct {
	Error string
}
