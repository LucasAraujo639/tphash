package diccionario

import (
	"fmt"
	"math"
)

type estado int

const (
	_VACIO estado = iota
	_OCUPADO
	_BORRADO
	_CAPACIDAD_INICIAL     = 21
	_AUMENTAR_CAPACIDAD    = 2.0
	_FACTOR_CARGA_DIMINUIR = 0.15
	_DISMINUIR_CAPACIDAD   = 0.5
	_INICIO                = 0
)

type campo[K comparable, V any] struct {
	clave  K
	valor  V
	estado estado
}

type hashCerrado[K comparable, V any] struct {
	tabla    []campo[K, V]
	cantidad int
	tam      int
	borrados int
}

func CrearHash[K comparable, V any]() Diccionario[K, V] {
	hash := new(hashCerrado[K, V])
	hash.tabla = crearTabla[K, V](_CAPACIDAD_INICIAL)
	hash.tam = _CAPACIDAD_INICIAL
	return hash
}
func crearTabla[K comparable, V any](capacidad int) []campo[K, V] {
	return make([]campo[K, V], capacidad)
}
func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func funcionHash[K comparable](clave K) int {
	claveBytes := convertirABytes(clave)
	h := 0
	for i := 0; i < len(claveBytes); i++ {
		h = h*31 + int(claveBytes[i])
	}
	return int(math.Abs(float64(h)))
}
func hashear[K comparable](clave K, capacidad int) int {
	return funcionHash(clave) % capacidad
}

// Guardar guarda el par clave-dato en el Diccionario. Si la clave ya se encontraba, se actualiza el dato asociado
func (hash *hashCerrado[K, V]) Guardar(clave K, valor V) {
	pos := hash.buscarPos(clave, hash.tam)
	if hash.tabla[pos].estado == _OCUPADO {
		hash.tabla[pos].valor = valor
		return
	}
	hash.tabla[pos].clave = clave
	hash.tabla[pos].valor = valor
	hash.tabla[pos].estado = _OCUPADO
	hash.cantidad++

	if int(hash.borrados+hash.cantidad) >= (hash.tam / 2) {
		hash.redimensionar(_AUMENTAR_CAPACIDAD)
	}
}

// busca la posicion de un elemento basada en una clave dada, devolviendo el índice correspondiente
func (hash *hashCerrado[K, V]) buscarPos(clave K, capacidad int) int {
	pos := hashear(clave, capacidad)
	for hash.tabla[pos].estado != _VACIO {
		if (hash.tabla[pos].estado == _OCUPADO) && (hash.tabla[pos].clave == clave) {
			break
		}

		pos++

		if pos >= capacidad {
			pos -= capacidad
		}
	}

	return pos
}

// Pertenece determina si una clave ya se encuentra en el diccionario, o no
func (hash *hashCerrado[K, V]) Pertenece(clave K) bool {
	return hash.tabla[hash.buscarPos(clave, hash.tam)].estado == _OCUPADO
}

// Obtener devuelve el dato asociado a una clave. Si la clave no pertenece, debe entrar en pánico con mensaje
// 'La clave no pertenece al diccionario'
func (hash *hashCerrado[K, V]) Obtener(clave K) V { //V
	pos := hash.buscarPos(clave, hash.tam)
	if hash.tabla[pos].estado == _OCUPADO {
		return hash.tabla[pos].valor
	}
	panic("La clave no pertenece al diccionario")
}

// // Borrar borra del Diccionario la clave indicada, devolviendo el dato que se encontraba asociado. Si la clave no
// // pertenece al diccionario, debe entrar en pánico con un mensaje 'La clave no pertenece al diccionario'
func (hash *hashCerrado[K, V]) Borrar(clave K) V {

	if float32(hash.borrados+hash.cantidad)/float32(hash.tam) <= _FACTOR_CARGA_DIMINUIR {
		hash.redimensionar(_DISMINUIR_CAPACIDAD)
	}
	pos := hash.buscarPos(clave, hash.tam)
	if hash.tabla[pos].estado != _OCUPADO {
		panic("La clave no pertenece al diccionario")
	}
	hash.tabla[pos].estado = _BORRADO
	hash.cantidad--
	hash.borrados++

	return hash.tabla[pos].valor
}

// crea un nuevo diccionario hash con la capacidad duplicada y reubica los elementos, devuelve un true si se redimensiono
func (hash *hashCerrado[K, V]) redimensionar(cap float64) {
	nuevaCapacidad := int(float64(hash.tam) * cap)
	nuevoCampo := crearTabla[K, V](nuevaCapacidad)
	hash.reubicarDatos(nuevoCampo, nuevaCapacidad)
}

// reubica todos los datos de mi tabla vieja a una nueva tabla hash redimensionada
func (hash *hashCerrado[K, V]) reubicarDatos(campo []campo[K, V], capacidad int) {
	aReubicar := hash.cantidad
	i := _INICIO

	for aReubicar > _INICIO {
		if hash.tabla[i].estado == _OCUPADO {
			guardar(campo, hash.tabla[i].clave, hash.tabla[i].valor, capacidad)
			aReubicar--
		}
		i++
	}

	hash.tabla = campo
	hash.borrados = 0
	hash.tam = capacidad

}

// funcion auxiliar para la redimension que guarda la clave el valor y el estado en mi nueva tabla hash
func guardar[K comparable, V any](campo []campo[K, V], clave K, valor V, capacidad int) {
	pos := buscarVacio(campo, clave, capacidad)
	campo[pos].clave = clave
	campo[pos].valor = valor
	campo[pos].estado = _OCUPADO
}

// busca una posicion vacia en donde ubicar la clave en la nueva tabla de hash
func buscarVacio[K comparable, V any](campo []campo[K, V], clave K, capacidad int) int {
	pos := hashear(clave, capacidad)
	for campo[pos].estado == _OCUPADO {
		pos++

		if pos >= capacidad {
			pos -= capacidad
		}
	}
	return pos
}

// Cantidad devuelve la cantidad de elementos dentro del diccionario
func (hash *hashCerrado[K, V]) Cantidad() int {
	return hash.cantidad
}

// Iterar itera internamente el diccionario, aplicando la función pasada por parámetro a todos los elementos del
// mismo
func (hash *hashCerrado[K, V]) Iterar(visitar func(clave K, valor V) bool) {
	for i := _INICIO; i < hash.tam; i++ {
		if hash.tabla[i].estado == _OCUPADO && !visitar(hash.tabla[i].clave, hash.tabla[i].valor) {
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
	return nuevoIter

}
func buscarIndiceProximaPosicion[K comparable, V any](iter iteradorDiccionario[K, V]) int {
	for i := iter.pos; i < iter.hashAsociado.tam; i++ {
		if iter.hashAsociado.tabla[i].estado == _OCUPADO {
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

// VerActual devuelve la clave y el dato del elemento actual en el que se encuentra posicionado el iterador.
// Si no HaySiguiente, debe entrar en pánico con el mensaje 'El iterador termino de iterar'
func (iter iteradorDiccionario[K, V]) VerActual() (K, V) {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	return iter.hashAsociado.tabla[iter.pos].clave, iter.hashAsociado.tabla[iter.pos].valor
}

// Siguiente si HaySiguiente avanza al siguiente elemento en el diccionario. Si no HaySiguiente, entonces debe
// entrar en pánico con mensaje 'El iterador termino de iterar'
func (iter *iteradorDiccionario[K, V]) Siguiente() {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	iter.pos++
	iter.pos = buscarIndiceProximaPosicion(*iter)

}
