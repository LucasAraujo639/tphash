package diccionario

import (
	"fmt"
	"math"
)

const (
	_CAPACIDAD_INICIAL  = 5
	_FACTOR_CARGA       = 0.7
	_AUMENTAR_CAPACIDAD = 2
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
	if hash.tabla[indice].estado == OCUPADO {
		hash.tabla[indice].valor = valor
		return false
	}
	hash.tabla[indice].clave = clave
	hash.tabla[indice].valor = valor
	hash.tabla[indice].estado = OCUPADO
	hash.tam++
	hash.cantidad++

	if int(hash.borrados+hash.cantidad) >= (hash.tam / 2) {
		redimensionar(hash)
	}
	return true
}
func hashing[K comparable](clave K, capacidad int) int {
	return funcionHash(clave) % capacidad
}

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
func guardar[K comparable, V any](campo []campo[K, V], clave K, valor V, capacidad int) {
	indice := buscarRedimension(campo, clave, capacidad)
	campo[indice].clave = clave
	campo[indice].valor = valor
	campo[indice].estado = OCUPADO
}
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
func redimensionar[K comparable, V any](hash *hashCerrado[K, V]) bool {
	nuevaCapacidad := hash.tam * 2
	nuevoCampo := make([]campo[K, V], nuevaCapacidad)
	fmt.Println("se redimensiono")
	reubicarDatos(hash, nuevoCampo, nuevaCapacidad)
	fmt.Println("se redimensiono")
	return true
}

// // Pertenece determina si una clave ya se encuentra en el diccionario, o no
func (hash hashCerrado[K, V]) Pertenece(clave K) bool {
	return hash.tabla[hash.buscar(clave, hash.tam)].estado == OCUPADO
}

// Obtener devuelve el dato asociado a una clave. Si la clave no pertenece, debe entrar en pánico con mensaje
// 'La clave no pertenece al diccionario'
func (hash hashCerrado[K, V]) Obtener(clave K) { //V
	// pos := funcionHash(clave, hash.tam)
	// for hash.tabla[pos].estado != VACIO{
	// 	if hash.tabla[pos].clave == clave{
	// 		return hash.tabla.valor
	// 	}
	// 	pos++
	// }
	// return panic("la clave no pertenece al diccionario")

}

// // Borrar borra del Diccionario la clave indicada, devolviendo el dato que se encontraba asociado. Si la clave no
// // pertenece al diccionario, debe entrar en pánico con un mensaje 'La clave no pertenece al diccionario'
func (hash *hashCerrado[K, V]) Borrar(clave string) string {

	return "a"
}

// Cantidad devuelve la cantidad de elementos dentro del diccionario
func (hash hashCerrado[k, v]) Cantidad() int {
	return hash.cantidad
}

// Iterar itera internamente el diccionario, aplicando la función pasada por parámetro a todos los elementos del
// // mismo
// func (hash hashCerrado[K, V]) Iterar(visitar func(clave K, valor V) bool) {
// 	pos := 0
// 	for hash.tabla[pos].estado != OCUPADO {
// 		pos++
// 	}

// 	for hash.tabla[pos] != nil &&  visitar(punteroIter.clave, punterIter.valor){
// 		pos++
// 	}
// }

// type iteradorDiccionario [K comparable, V any] struct{
// 	pos int
// 	siguiente int
// 	diccionarioAsociado *hashCerrado[K, V]
// }
// 	// Iterador devuelve un IterDiccionario para este Diccionario
// 	func (iter iterDiccionario[K, V])Iterador() IterDiccionario[K, V] {
// 	nuevoIter := new(IteradorDiccionario[K, V])
// 	nuevoIter.pos = buscarPrimeraPosicion(nuevoIter.diccionarioAsociado, 0)
// 	nuevoIter.siguiente = pos+1
// 	return nuevoIter

// }
// func avanzarAlSiguienteOcupado(hash hashCerrado[K, V], posActual int) int{

// 	for hash.tabla[posActual].estado != OCUPADO {
// 		posActual++
// 	}
// 	return posActual
// }

// 	// HaySiguiente devuelve si hay más datos para ver. Esto es, si en el lugar donde se encuentra parado
// 	// el iterador hay un elemento.
// func  (iter iterDiccionario[K, V])()HaySiguiente() bool {
// 	return iter.pos != iter.hashCerrado.tam
// }

// 	// VerActual devuelve la clave y el dato del elemento actual en el que se encuentra posicionado el iterador.
// 	// Si no HaySiguiente, debe entrar en pánico con el mensaje 'El iterador termino de iterar'
// 	func (iter iterDiccionario[K,V]) VerActual() (K, V){
// 		if !iter.HaySiguiente(){
// 			panic("el iterador termino de iterar")
// 		return iter.hashCerrado.tabla[iter.pos].clave, iter.hashCerrado.tabla[iter.pos].valor
// 	}

// 	// Siguiente si HaySiguiente avanza al siguiente elemento en el diccionario. Si no HaySiguiente, entonces debe
// 	// entrar en pánico con mensaje 'El iterador termino de iterar'
// 	func (iter iterDiccionario[K,V]) Siguiente(){
// 		iter.pos++
// 		if !iter.HaySiguiente(){
// 			panic("el iterador termino de iterar")
// 			}
// 			for iter.hashCerrado.tabla[iter.pos].estado != OCUPADO {
// 			posActual++
// 		}

// 	}
