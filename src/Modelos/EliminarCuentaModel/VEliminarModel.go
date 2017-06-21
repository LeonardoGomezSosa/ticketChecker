package EliminarModel

//SEliminarVista estructura Auxiliar para enviar listas a la vista
type SEliminarVista struct {
	ID      string
	SEstado bool
	SMsj    string
	Error   []Errores
}

//Errores interfaz para capturar errores
type Errores struct {
	Error string
}
