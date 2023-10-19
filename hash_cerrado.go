package diccionario

import (
	"fmt"
	"math"
)

const (
	_CAPACIDAD_INICIAL  = 10
	_FACTOR_CARGA       = 0.7
	_AUMENTAR_CAPACIDAD = 2
	_INICIO             = 0
)

type Estado int

const (
	VACIO Estado = iota
	OCUPADO
	BORRADO
)

type campo[K comparable, V any] struct {
	clave  K
	valor  V
	estado Estado
}

type hashCerrado[K comparable, V any] struct {
	tabla    []campo[K, V]
	cantidad int
	tam      int
	borrados int
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	hash := new(hashCerrado[K, V])
	hash = &hashCerrado[K, V]{
		tabla: make([]campo[K, V], _CAPACIDAD_INICIAL),
		tam:   _CAPACIDAD_INICIAL,
	}
	return hash
}
func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func funcionHash[K comparable](clave K) int {
	claveBytes := convertirABytes(clave)
	p := 0
	for i := 0; i < len(claveBytes); i++ {
		p = p*31 + int(claveBytes[i])
	}
	return int(math.Abs(float64(p)))
}

// Guardar guarda el par clave-dato en el Diccionario. Si la clave ya se encontraba, se actualiza el dato asociado
func (hash *hashCerrado[K, V]) Guardar(clave K, valor V) bool {
	indice := hash.buscar(clave, hash.tam)
	fmt.Println("i", indice)
	if hash.tabla[indice].estado == OCUPADO {
		hash.tabla[indice].valor = valor
		return false
	}
	hash.tabla[indice].clave = clave
	hash.tabla[indice].valor = valor
	hash.tabla[indice].estado = OCUPADO
	hash.cantidad++

	if int(hash.borrados+hash.cantidad) >= (hash.tam / 2) {
		redimensionar(hash)
	}
	return true
}
func hashing[K comparable](clave K, capacidad int) int {
	return funcionHash(clave) % capacidad
}

// busca la posicion de un elemento basada en una clave dada, devolviendo el índice correspondiente
func (hash hashCerrado[K, V]) buscar(clave K, capacidad int) int {
	indice := hashing(clave, capacidad)
	for hash.tabla[indice].estado != VACIO {
		if (hash.tabla[indice].estado == OCUPADO) && (hash.tabla[indice].clave == clave) {
			break
		}

		indice++

		if indice >= capacidad {
			indice -= capacidad
		}
	}

	return indice
}

// reubica todos los datos de mi tabla vieja a una nueva tabla hash redimensionada
func reubicarDatos[K comparable, V any](hash *hashCerrado[K, V], campo []campo[K, V], capacidad int) {
	aVisitar := hash.cantidad

	for i := 0; aVisitar > 0; i++ {
		if hash.tabla[i].estado == OCUPADO {
			guardar(campo, hash.tabla[i].clave, hash.tabla[i].valor, capacidad)
			aVisitar--
		}
	}

	hash.tabla = campo
	hash.tam = capacidad
	hash.borrados = 0
}

// funcion auxiliar para la redimension que guarda la clave el valor y el estado en mi nueva tabla hash
func guardar[K comparable, V any](campo []campo[K, V], clave K, valor V, capacidad int) {
	indice := buscarRedimension(campo, clave, capacidad)
	campo[indice].clave = clave
	campo[indice].valor = valor
	campo[indice].estado = OCUPADO
}

// busca una posicion vacia en donde ubicar la clave en la nueva tabla de hash
func buscarRedimension[K comparable, V any](campo []campo[K, V], clave K, capacidad int) int {
	indice := hashing(clave, capacidad)
	for campo[indice].estado == OCUPADO {
		indice++
		if indice >= capacidad {
			indice -= capacidad
		}
	}
	return indice
}

// crea un nuevo diccionario hash con la capacidad duplicada y reubica los elementos, devuelve un true si se redimensiono
func redimensionar[K comparable, V any](hash *hashCerrado[K, V]) bool {
	nuevaCapacidad := hash.tam * 2
	nuevoCampo := make([]campo[K, V], nuevaCapacidad)

	reubicarDatos(hash, nuevoCampo, nuevaCapacidad)
	fmt.Println("se redimensiono")
	fmt.Println("tamaño nuevo", hash.tam)
	return true
}

// Pertenece determina si una clave ya se encuentra en el diccionario, o no
func (hash hashCerrado[K, V]) Pertenece(clave K) bool {
	return hash.tabla[hash.buscar(clave, hash.tam)].estado == OCUPADO
}

// Obtener devuelve el dato asociado a una clave. Si la clave no pertenece, debe entrar en pánico con mensaje
// 'La clave no pertenece al diccionario'
func (hash hashCerrado[K, V]) Obtener(clave K) V { //V
	indice := hash.buscar(clave, hash.tam)
	if hash.tabla[indice].estado == OCUPADO {
		return hash.tabla[indice].valor
	}
	panic("la clave no pertenece al diccionario")
}

// // Borrar borra del Diccionario la clave indicada, devolviendo el dato que se encontraba asociado. Si la clave no
// // pertenece al diccionario, debe entrar en pánico con un mensaje 'La clave no pertenece al diccionario'
func (hash *hashCerrado[K, V]) Borrar(clave K) V {

	indice := hash.buscar(clave, hash.tam)
	if hash.tabla[indice].estado != OCUPADO {
		panic("La clave no pertenece al diccionario")
	}
	hash.tabla[indice].estado = BORRADO
	hash.cantidad--
	hash.borrados++
	return hash.tabla[indice].valor

}

// Cantidad devuelve la cantidad de elementos dentro del diccionario
func (hash hashCerrado[K, V]) Cantidad() int {
	return hash.cantidad
}

// Iterar itera internamente el diccionario, aplicando la función pasada por parámetro a todos los elementos del
// mismo
func (hash hashCerrado[K, V]) Iterar(visitar func(clave K, valor V) bool) {
	for i := 0; i < hash.tam; i++ {
		if hash.tabla[i].estado == OCUPADO && !visitar(hash.tabla[i].clave, hash.tabla[i].valor) {
			break
		}
	}
}

type iteradorDiccionario[K comparable, V any] struct {
	pos          int
	hashAsociado *hashCerrado[K, V]
}

// Iterador devuelve un IterDiccionario para este Diccionario
func (hash *hashCerrado[K, V]) Iterador() IterDiccionario[K, V] {
	nuevoIter := new(iteradorDiccionario[K, V])
	nuevoIter.hashAsociado = hash
	nuevoIter.pos = buscarIndiceProximaPosicion(*nuevoIter)
	fmt.Println("creo hash posicion en la encuentra", nuevoIter.pos)
	return nuevoIter

}
func buscarIndiceProximaPosicion[K comparable, V any](iter iteradorDiccionario[K, V]) int {
	fmt.Println("buscar pos", iter.hashAsociado.tam)
	for i := iter.pos; i < iter.hashAsociado.tam; i++ {
		if iter.hashAsociado.tabla[i].estado == OCUPADO {
			return i
		}
	}
	return iter.hashAsociado.tam
}

// HaySiguiente devuelve si hay más datos para ver. Esto es, si en el lugar donde se encuentra parado
// el iterador hay un elemento.
func (iter iteradorDiccionario[K, V]) HaySiguiente() bool {
	return iter.pos != iter.hashAsociado.tam
}

// 	// VerActual devuelve la clave y el dato del elemento actual en el que se encuentra posicionado el iterador.
// 	// Si no HaySiguiente, debe entrar en pánico con el mensaje 'El iterador termino de iterar'
// 	func (iter iterDiccionario[K,V]) VerActual() (K, V){
// 		if !iter.HaySiguiente(){
// 			panic("el iterador termino de iterar")
// 		return iter.hashCerrado.tabla[iter.pos].clave, iter.hashCerrado.tabla[iter.pos].valor
// 	}

// Siguiente si HaySiguiente avanza al siguiente elemento en el diccionario. Si no HaySiguiente, entonces debe
// entrar en pánico con mensaje 'El iterador termino de iterar'
func (iter *iteradorDiccionario[K, V]) Siguiente() {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	iter.pos++
	iter.pos = buscarIndiceProximaPosicion(*iter)
	fmt.Println("posicion siguiente", iter.pos)

}
