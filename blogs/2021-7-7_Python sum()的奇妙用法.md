```python
a = [1, 4, ['a', 'a'], 5, 3]
b = sum([[i] if isinstance(i, int) else i for i in a], [])
print(b)
```

> [1, 4, 'a', 'a', 5, 3]
