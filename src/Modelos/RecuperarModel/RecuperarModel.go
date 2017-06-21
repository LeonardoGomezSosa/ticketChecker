package RecuperarModel

//SRecuperarVista estructura Auxiliar para enviar listas a la vista
type SRecuperarVista struct {
	ID      string
	SEstado bool
	SMsj    string
	Error   []ErroresR
}

//Errores interfaz para capturar errores
type ErroresR struct {
	Error string
}
