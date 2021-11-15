# Python defaultdict的特殊应用

v2上看到的，记录一下。

[最近发现 defaultdict 的一个奇技淫巧](https://www.v2ex.com/t/759173)

> defaultdict(default_factory[, ...]) --> dict with default factoryThe default factory is called without arguments to producea new value when a key is not present, in __getitem__ only.A defaultdict compares equal to a dict with the same items.All remaining arguments are treated the same as if they werepassed to the dict constructor, including keyword arguments.

可以用lambda让defaultdict生成一个去重的自增索引

```python
from collections import defaultdict

ind = defaultdict(lambda: len(ind))
var = ind["test_a"]
var1 = ind["test_b"]
var2 = ind["test_a"]
print(var, var1, var2)
```

> 0 1 0