package diccionario

type funcCmp[K comparable] func(K, K) int

type nodoAbb[K comparable, V any] struct {
	izquierdo *nodoAbb[K, V]
	derecho   *nodoAbb[K, V]
	clave     K
	dato      V
}

type abb[K comparable, V any] struct {
	raiz     *nodoAbb[K, V]
	cantidad int
	cmp      funcCmp[K]
}

// type iterDiccionarioOrdenado[K comparable, V any]{
// 	dato
// 	desde clave
// 	hasta clave
// 	abb *abb[K, V]
// }

// Guardar guarda el par clave-dato en el Diccionario. Si la clave ya se encontraba, se actualiza el dato asociado
func Guardar(clave K, dato V) {

}

// // Pertenece determina si una clave ya se encuentra en el diccionario, o no
// func Pertenece(clave K) bool
// // Obtener devuelve el dato asociado a una clave. Si la clave no pertenece, debe entrar en pánico con mensaje 'La clave no pertenece al diccionario'
// func Obtener(clave K) V{

// }

// // Borrar borra del Diccionario la clave indicada, devolviendo el dato que se encontraba asociado. Si la clave no pertenece al diccionario, debe entrar en pánico con un mensaje 'La clave no pertenece al diccionario'
// func Borrar(clave K) V{

// }
// // Cantidad devuelve la cantidad de elementos dentro del diccionario
// func Cantidad() int{

// }
// // Iterar itera internamente el diccionario, aplicando la función pasada por parámetro a todos los elementos del mismo
// func Iterar(func(clave K, dato V) bool){

// }
// // Iterador devuelve un IterDiccionario para este Diccionario
// func Iterador() IterDiccionario[K, V]{

// }

// 	// IterarRango itera sólo incluyendo a los elementos que se encuentren comprendidos en el rango indicado, incluyéndolos en caso de encontrarse
// 	func IterarRango(desde *K, hasta *K, visitar func(clave K, dato V) bool){

// 	}

// 	// Iterador Rango crea un IterDiccionario que sólo itere por las claves que se encuentren en el rango indicado
// 	func IteradorRango(desde *K, hasta *K) IterDiccionario[K, V]{

// 	}

// //------ iter externo -----

// // HaySiguiente devuelve si hay más datos para ver. Esto es, si en el lugar donde se encuentra parado el iterador hay un elemento.
// func HaySiguiente() bool{

// }

// // VerActual devuelve la clave y el dato del elemento actual en el que se encuentra posicionado el iterador. Si no HaySiguiente, debe entrar en pánico con el mensaje 'El iterador termino de iterar'
// 	func VerActual() (K, V){

// 	}

// // Siguiente si HaySiguiente, devuelve la clave actual (equivalente a VerActual, pero únicamente la clave), y además avanza al siguiente elemento en el diccionario. Si no HaySiguiente, entonces debe entrar en pánico con mensaje 'El iterador termino de iterar'
// func Siguiente() K{

// 	}
