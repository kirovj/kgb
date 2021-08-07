# 用 Java 重构 Python 项目的感概

> 虽然还有许多工作需要做，但无论怎样还是终于接近尾声了。在这里也记录一下这两个月重构过程中遇到的困难和感慨吧。

## 看大型项目还是很吃力

尽管工程结构不是很复杂，逻辑与设计模式比较清晰（毕竟是 python ）。但是整体看下来还是比较头大，需要反复 debug 才能让整个运行逻辑比较清晰。

最大的阻力自然是类型。虽然项目中大部分都有辅助类型推断（ PEP8 ），但是吧比如：

```python
from typing import List
def foo() -> List[Union(A, B)]:
    pass
```

这在 Java 中你不得不加以抽象或者其他方式来解决。

cpython 的特性也是一大阻力，因为底层是c， IDE 无法直接索引到c标准库源码，所以你会发现很多标准库的方法点进去只有一个 pass （其实这个也是 IDE 帮你显示的，因为底层是c模块）。

当然了，还有 python 一大堆的高级语法特性。比如 itertools 包中的许多方法（ groupby ...）、列表生成式中的列表生成式还嵌套了 map | filter 、闭包、无处不在的 yield | yield from ...

---

## PDF - 给计算机看的文字

重构的项目是和 pdf 有关的。『 PDF - 给计算机看的文字』这句话简单来说， *_PDF 就不是被设计出来编辑的，而是实体文档在数字世界中的影像。_*

先看一下 word 底层的源码：
```
<w:pPr>
    <w:jc w:val="center"/>
    <w:rPr><w:rFonts w:ascii="Times" w:hAnsi="Times"/><w:lang w:val="en-US"/></w:rPr>
</w:pPr>
<w:r w:rsidRPr="003C75CF">
    <w:rPr><w:rFonts w:ascii="Times" w:hAnsi="Times"/></w:rPr>
    <w:t>Hello world!</w:t>
</w:r>
```
> 这就是 docx 格式这类标记语言 （ Markup Language ）文档的特征：在纯文本上包裹各种「标签」（ tag ）来描述文本的样式（颜色、位置、字体等等），从而获得格式丰富多变的文档。常见的网页（ HTML ）、 Evernote 笔记（ ENML ）所用的语法本质上都是标记语言，区别只在于支持的标签各不相同、因此能实现的格式有多有少罢了。

看一下 pdf 的源码：

```
BT
    1 0 0 1 1036 572 Tm
    /TT1 12 Tf
    [ (He) 24 (l) -48 (l) -48 (o) ] TJ
ET
BT
    1 0 0 1 1147 572 Tm
    /TT1 12 Tf
    ( ) Tj
ET
BT
    1 0 0 1 1160 572 Tm
    /TT1 12 Tf
    [ (w) 24 (or) -84 (l) -24 (d) ] TJ
ET
```
> 以上是 PDF 的部分源码，显示了字符串 Hello world

翻译一下：
```
【文字开始】
    缩放比例1倍 坐标(1036,572) 【文字定位】
    /TT1 12磅 【选择字体】
    [ (He) 间距24 (l) 间距-48 (l) 间距-48 (o) ] 【绘制文字】
【文字结束】
【文字开始】
    缩放比例1倍 坐标(1147,572) 【文字定位】
    /TT1 12磅 【选择字体】
    (空格) 【绘制文字】
【文字结束】
【文字开始】
    缩放比例1倍 坐标(1060,572) 【文字定位】
    /TT1 12磅 【选择字体】
    [ (w) 间距24 (or) 间距-84 (l) 间距-24 (d) ] 【绘制文字】
【文字结束】
```
>  不同于 word ， pdf 的语言则更像是在控制机器，定位、调整、落笔、抬起，移动到下一行；如此重复。

当然，以上只是 pdf 中很小部分的介绍，详细的可以看下 adobe 公司的[ pdf 标准文件](https://www.adobe.com/content/dam/acom/en/devnet/pdf/pdfs/PDF32000_2008.pdf)。

以及一些延伸阅读：

* [PDF Reference and Adobe Extensions to the PDF Specification](https://www.adobe.com/devnet/pdf/pdf_reference.html): Adobe 的 PDF 文档专页，提供了 PDF 1.7（成为 ISO 标准的版本）和 Adobe 基于 PDF 1.7 版所做功能扩展的文档。
* [Understanding the PDF File Format](https://blog.idrsolutions.com/2013/01/understanding-the-pdf-file-format-overview): 帮助理解 PDF 格式的详尽系列文章，每篇都不长，语言比较浅显易懂。
* [Planet PDF](http://www.planetpdf.com/): Foxit 旗下的 PDF 专题网站。页面设计很老旧，但是能找到不少探究 PDF 格式细节的资料。 
* [John Whitington, PDF Explained (O'Reilly Media, 2011)](https://github.com/zxyle/PDF-Explained): 动物园丛书里介绍 PDF 的不止一本，这本年代偏久远，但是个人认为是相比之下讲得比较清晰的。
* [GitHub - PDF 101 - Learn and Play with PDF Source Code](https://github.com/angea/PDF101): 一个非常有意思的 GitHub 仓库，里面提供了大量手敲代码制成的 PDF，用文本编辑器打开按照里面的注释就可以像玩游戏一样做各种「实验」，对理解 PDF 格式非常有帮助。

> 👆🏻感谢：[理解数字世界中的纸张：PDF | 科普](https://zhuanlan.zhihu.com/p/44360779)

---

## JNI 与 SWIG 

todo

---

## 算法的不足

* [维诺图（ Voronoi Diagram ）](https://zh.wikipedia.org/wiki/%E6%B2%83%E7%BD%97%E8%AF%BA%E4%BC%8A%E5%9B%BE)

* [连通分量/图（ Connected Component ）](https://zh.wikipedia.org/wiki/%E8%BF%9E%E9%80%9A%E5%9B%BE)