v0.0.5 - до сюда пример про яйцо  
v0.0.13 - splitFlexCol первый пошел.

v.0.0.14 - меняется версия Gio (go.mod), все нахрен слетает. все исправляем по надобности. (исправлям op.Save)  
splitFlexCol - работает продолжаем

```
[изучаем тут](https://git.sr.ht/~eliasnaur/gio/commit/6e9574245074656fe272928ec0523be4749a80a9)

API change: replace op.Save/Load with explicit Push/Pop scopes for
op.TransformOps, pointer.AreaOps, clip.Ops.

Before this change, op.Save and Load was used to save and restore the
state:

    ops := new(op.Ops)
    // Save state.
    state := op.Save(ops)
    // Apply offset.
    op.Offset(...).Add(ops)
    // Draw with offset applied.
    draw(ops)
    // Restore state.
    state.Load()

The example above now becomes:

    ops := new(op.Ops)
    // Push offset to the transformation stack.
    stack := op.Offset(...).Push(ops)
    // Draw with offset applied.
    draw(ops)
    // Restore state.
    stack.Pop()
```
