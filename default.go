package stackable

type Default interface {
    Default() any
}

func DefaultValue[T Default]() T {
    var t T
    return t.Default().(T)
}
