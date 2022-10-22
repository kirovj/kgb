> 学习语言的最好方法就是用它码代码，而爬虫，则再合适不过

## 错过的价格一去不复返

19年初接触 csgo，经历了白给->练枪->学道具->跑图->看职业选手pov，也见证了饰品价格的一路飞涨。

从一千到六千的蝴蝶刀，从几块钱的贴纸到现在几万一张有价无市，几十万的玄学661、全蓝淬火。这个市场是畸形的，无论是极低的开箱出金概率，还是诈骗、屯货的奸商、垄断的交易平台、以及没有人想要把自己买来的饰品降价出手。

每当普通玩家想要购买一把刀、一双手套，或者是好看的枪皮，都会考虑许久。买了之后价格会不会跌？卖的时候能不能至少卖出原价？而当你想要看看价格历史走势的时候，却发现垄断的平台只让你看到六个月的数据──这还需要积分换取。

这个需求就是由此而来，如果一切顺利的话，我想要在几年后看到这段时间饰品价格的起起落落甚至一窥游戏的兴衰。

## Rust, Reqwest, Serde, Sqlite

语言，http，序列化/反序列化，数据库。

学习语言的最好方法就是用它码代码，而爬虫，则再合适不过。从 python 到 java 再到 go 最后是 rust，我记得学习他们的第一个项目都是爬虫。

可以是几行代码的脚本，也可以是几个函数，面向对象？当然可以，工程化？你想要也可以。由浅至深，你可以了解最基本的 http 协议，html 解析，正则表达式，json 反序列化，数据库，并发，反反爬，逆向，验证码破解……

## http: Reqwest

http.rs
```rust
use std::collections::HashMap;
use reqwest::header::{HeaderMap, HeaderValue};
use crate::utils::UA;

fn make_headers() -> HeaderMap {
    let mut map = HeaderMap::new();
    map.insert("user-agent", HeaderValue::from_static(UA));
    map
}

lazy_static! {
    static ref CLIENT: reqwest::blocking::Client = reqwest::blocking::Client::builder()
        .default_headers(make_headers())
        .build()
        .unwrap();
}

enum Method {
    Get,
    PostJson,
}

fn request(
    url: &str,
    method: Method,
    data: Option<&HashMap<&str, &str>>,
) -> Result<String, reqwest::Error> {
    let mut times = 1;
    loop {
        // 失败重试的简陋实现
        let r = match method {
            Method::Get => CLIENT.get(url).send()?.text(),
            Method::PostJson => CLIENT.post(url).json(data.unwrap()).send()?.text(),
        };
        if !r.is_ok() && times < 3 {
            times += 1;
            continue;
        }
        break r;
    }
}
```

## database: Rusqlite

db.rs

```rust
use crate::item::{Item, PriceInfo};
use rusqlite::{params, Connection, Result};
use std::path::Path;

fn create_db(db_file: &str) -> Result<(), rusqlite::Error> {
    let conn = Connection::open(db_file)?;
    let _ = conn.execute(
        "create table Item (
        id        INTEGER PRIMARY KEY AUTOINCREMENT,
        name      VARCHAR NOT NULL,
        class     CHAR(16) NOT NULL,
        typo      CHAR(32) NOT NULL,
        ware      CHAR(16) NOT NULL,
        quality   CHAR(16) NOT NULL,
        rarity    CHAR(16) NOT NULL,
        stat_trak INTEGER NOT NULL
    )",
        [],
    );
    let _ = conn.execute(
        "create table PriceInfo (
        id      INTEGER PRIMARY KEY AUTOINCREMENT,
        item_id INTEGER NOT NULL,
        date    CHAR(16) NOT NULL,
        price   NUMERIC(10,1) NOT NULL
    )",
        [],
    );
    Ok(())
}

pub struct DbHelper {
    conn: Connection,
}

#[allow(unused)]
impl DbHelper {
    pub fn new(db_file: &str) -> DbHelper {
        if !Path::new(db_file).exists() {
            create_db(db_file);
        }
        let conn = Connection::open(db_file).unwrap();
        DbHelper { conn }
    }

    fn find_item_id(&self, item: &Item) -> Result<u32> {
        let mut stmt = self
            .conn
            .prepare(
                "SELECT id from Item
            where name = ?1
            and class = ?2
            and typo = ?3
            and ware = ?4
            and quality = ?5
            and rarity = ?6
            and stat_trak = ?7",
            )
            .unwrap();
        stmt.query_row(
            params![
                item.name,
                item.class,
                item.typo,
                item.ware,
                item.quality,
                item.rarity,
                item.stat_trak
            ],
            |row| row.get(0),
        )
    }

    pub fn get_item_id(&self, item: &Item) -> Option<u32> {
        match self.find_item_id(item) {
            Ok(id) => Some(id),
            _ => {
                self.add_item(item);
                match self.find_item_id(item) {
                    Ok(_id) => Some(_id),
                    _ => None,
                }
            }
        }
    }

    pub fn add_item(&self, item: &Item) {
        self.conn.execute(
            "INSERT INTO Item
            (name, class, typo, ware, quality, rarity, stat_trak)
            VALUES(?1, ?2, ?3, ?4, ?5, ?6, ?7)",
            params![
                item.name,
                item.class,
                item.typo,
                item.ware,
                item.quality,
                item.rarity,
                item.stat_trak
            ],
        );
    }

    pub fn find_price_info_id(&self, price_info: &PriceInfo) -> Result<usize> {
        let mut stmt = self
            .conn
            .prepare("SELECT id from PriceInfo where item_id = ?1 and date = ?2")
            .unwrap();
        stmt.query_row(params![price_info.item_id, price_info.date], |row| {
            row.get(0)
        })
    }

    pub fn add_price_info(&self, price_info: &PriceInfo) {
        match self.find_price_info_id(price_info) {
            Ok(id) => {
                self.conn.execute(
                    "update PriceInfo set price = ?1 where id = ?2",
                    params![price_info.price, id],
                );
            }
            _ => {
                self.conn.execute(
                    "INSERT INTO PriceInfo (item_id, date, price) VALUES(?1, ?2, ?3)",
                    params![price_info.item_id, price_info.date, price_info.price,],
                );
            }
        }
    }
}
```

## crawler

crawler.rs

```rust
pub trait Crawl {
    fn name(&self) -> &str;

    fn alert(&self, message: &str) {
        let message = format!("[{}] {}: {}", utils::current_time(), self.name(), message);
        println!("{}", message);
        utils::alert(message.as_str());
    }

    fn db(&self) -> &DbHelper;

    fn run(&self);

    fn persistent(&self, item: Item, price: &str) {
        match price.parse::<f32>() {
            Ok(p) => match self.db().get_item_id(&item) {
                None => {}
                Some(id) => {
                    self.db().add_price_info(&PriceInfo::new(
                        id,
                        utils::current_date(),
                        utils::round(p),
                    ));
                }
            },
            _ => self.alert(format!("parse price {} err", price).as_str()),
        }
    }

    fn sleep(&self) {
        thread::sleep(time::Duration::from_secs(
            rand::thread_rng().gen_range(20..40),
        ));
    }
}
```

## License

    GPL-2.0

## Warning

爬虫写得好──牢饭吃得早，切记。所以：

```
Open source if you want to create your own version.
如果你想要开发自己的版本，请开源并遵守协议。

Do not use it in an illegal way.
不要非法使用。

Do not use data in an illegal way.
不要非法使用通过它爬的数据。

I take no responsibility.
我不负任何责任。
```

> Github: [github.com/kirovj/csgo-crawler](https://github.com/kirovj/rscsgo)