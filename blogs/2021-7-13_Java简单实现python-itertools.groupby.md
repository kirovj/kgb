# Java简单实现python-itertools.groupby()

```python
class groupby(object):
    """
    make an iterator that returns consecutive keys and groups from the iterable
    
      iterable
        Elements to divide into groups according to the key function.
      key
        A function for computing the group category for each element.
        If the key function is not specified or is None, the element itself
        is used for grouping.
    """
```

### _python usage:_

```python
from itertools import groupby

data = [(False, 1), (False, 5), (True, 3), (True, 1), (False, 1), ]

for tag, d in groupby(data, key=lambda x: x[0]):
    print(tag, list(d))
```

False [(False, 1), (False, 5)]  
True  [(True, 3), (True, 1)]  
False [(False, 1)]

> 简单来说， 就是给一个列表，按照该列表元素的某个特征进行分组但是并不聚合

---

### _Java实现_

首先，定义一个类表示组
```java
public static class Group<T, E> {
    // 元素的分组特征
	final T type;
    // 要分组的 List
	List<E> data;

	public Group(T type) {
		this.type = type;
        this.data = new ArrayList();
	}

	public T type() {
		return type;
	}

	public List<E> data() {
		return data;
	}

	public void add(E e) {
		this.data.add(e);
	}
}
```

定义一个接口，包含两个方法：  
1. 通过一个元素获取它的特征
2. 比较两个特征是否一致

```java
public interface TypeDeal<T, E> {
    /**
     * 比较两个特征是否一致
     *
     * @param t1 last T
     * @param t2 next T
     * @param e next E
     * @return if e the same type
     */
    default boolean check(T t1, T t2, E e) {
        return t1.equals(t2);
    }

    /**
     * 通过一个元素获取它的特征
     * @param e item
     * @return item's type
     */
    T make(E e);
}
```

最后实现分组即可

```java
public <T, E> List<GroupBy<T, E>> build(List<E> items, TypeDeal<T, E> dealer) {
    List<GroupBy<T, E>> result = new ArrayList<>();
    if (items != null && !items.isEmpty()) {
        T t0 = dealer.make(items.get(0));
        Group<T, E> group = new Group<>(t0);
        for (E item : items) {
            T t1 =  dealer.make(item);
            if (!dealer.check(t0, t1){
                result.add(group);
                t0 = t1;
                group = new Group<>(t0);
            }
            group.add(item);
        }
        result.add(group);
    }
    return result;
}
```