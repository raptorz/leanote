# Pearlnote

[![Build Status](https://travis-ci.org/pearlnote/pearlnote.svg)](https://travis-ci.org/pearlnote/pearlnote)
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/pearlnote/pearlnote?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## 1. Introduction

Pearlnote, not just a notepad!
**Highlighted Features**

* Note-taking made easy: Pearlnote incorporates a clean and intuitive interface, the `tinymce` rich-text editor and a dedicated *markdown* editor, making your writing/typing more efficient and enjoyable. For more advanced users, we even offer `Vim` and `Emacs` writing modes to help boost your writing speed to another level.
* Knowledge management: The flexible and versatile notebook-note-tagging system of Pearlnote makes it an ideal tool for knowledge management.
* Sharing: Share your knowledge, thoughts and experiences with friends via Pearlnote. Invite your friends to join your notepad in the cloud.
* Cooperating: Collaborate with colleagues to improve skills, fertilize ideas and brainstorm on the fly.
* Blogging: Publish your work and make Pearlnote your personal blog.

**Other Features**

* Markdown syntax support
* Distraction-free writing mode
* `Vim` and `Emacs` editing mode
* Export notes to PDFs
* Batch note operation
* Customizable themes for blogging

## 2. Why we create Pearlnote

We have been using the popular note-taking software/service `Evernote` as our knowledge management tool on a daily basis. Benefited from and inspired by `Evernote`, we decided to create a brand-new tool that provides everything `Evernote` has to offer, plus a bunch of new features that `Evernote` failed to deliver, such as:

* A more powerful editor: `Evernote`'s editor lacks the functionalities of **document navigation**, **syntax based code rendering** (as a programmer, syntax highlighted code rendering is a necessity), **image resizing** and so forth.
* Everybody loves *markdown*, however `Evernote` simply wouldn't add it despite of years' of requests from users. So we will do the favor and bring a *markdown* enabled editor to you, guess what, it is also rendered in real-time!
* If you a developer and miss the `Vim` or `Emacs` ways of writing, we offer you the choice of `Vim` and `Emacs` editing modes. Equipped with *markdown* syntax for text formatting, you will never need to touch your mouse while writing.
* We love managing knowledge and thoughts as much as sharing them, so everybody has their own note account (`Evernote`, `Onenote`, `Google doc`, `Wiz note` etc.) and social media account (`Facebook`, `Wordpress`, blogs, etc.). But why can’t those two be one? Pearlnote makes this first step to bridge the private note-taking and public knowledge sharing seamlessly.
* A complete and all-platform (sorry Windows phone) covering software suite: that includes Pearlnote Web & Server (this repository), [Desktop app](https://github.com/pearlnote/desktop-app), [iOS](https://github.com/pearlnote/pearlnote-ios), [Android](https://github.com/pearlnote/pearlnote-android). And they are all open source!
* ......

## 3. How to get Pearlnote

The Pearlnote software suite contains: Pearlnote Web & Server (this repository), [Desktop app](https://github.com/pearlnote/desktop-app), [iOS](https://github.com/pearlnote/pearlnote-ios), [Android](https://github.com/pearlnote/pearlnote-android).

Interested in our product and want to try it out from your web browser? Welcome to sign up on https://pearlnote.com.

Feeling suspicious about how those note-taking companies treat your personal data? You can install Pearlnote on your server, and use Pearlnote App (Desktop, iOS, Android) to sync notes with your self-hosted server.

**Database abstraction and migration**

Pearlnote now supports both **MongoDB** and **PostgreSQL** databases with:
- Unified database interface for seamless database switching
- Bidirectional data migration tools (MongoDB ↔ PostgreSQL)
- Automatic ID mapping and data validation
- Configurable via `conf/app.conf`

See [DATABASE_ABSTRACTION_README.md](DATABASE_ABSTRACTION_README.md) for details.

More information about how to install Pearlnote please see:

* Pearlnote binary installation tutorial:
    * [Windows](https://github.com/pearlnote/pearlnote/wiki/pearlnote-source-installation-on-Windows-(En))
    * [Mac and Linux](https://github.com/pearlnote/pearlnote/wiki/pearlnote-binary-installation-on-Mac-and-Linux-(En))
* Pearlnote source installation tutorial:
    <!-- * [Windows](https://github.com/pearlnote/pearlnote/wiki/pearlnote-source-installation-on-Windows-(En)) -->
    * [Mac and Linux](https://github.com/pearlnote/pearlnote/wiki/Pearlnote-source-installation-on-Mac-and-Linux-(En))

## 4. Documentation

Please see [wiki](https://github.com/pearlnote/pearlnote/wiki) for detailed instruction on how to install Pearlnote on various platforms, trouble shooting and configuration explanations.

**Database Abstraction & Migration Documentation**:
- [DATABASE_ABSTRACTION_README.md](DATABASE_ABSTRACTION_README.md) - Project overview and quick start
- [docs/DATABASE_ABSTRACTION_GUIDE.md](docs/DATABASE_ABSTRACTION_GUIDE.md) - Complete usage guide

### Rename compatibility

Pearlnote keeps the existing HTTP API routes, request parameters, response fields, MongoDB BSON fields, and PostgreSQL columns unchanged. Existing deployments may continue pointing `db.dbname` at their current database; the new default database and deployment identifiers are `pearlnote`. The legacy `app.secretLeanote` setting remains supported as a fallback for upgrades.


## 5. How to develop Pearlnote

If you are a developer yourself and feel like to build on top of Pearlnote, please refer to [How-to-develop-pearlnote](https://github.com/pearlnote/pearlnote/wiki/How-to-develop-pearlnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91pearlnote).


## 6. Contributions

Like or dislike Pearlnote, please leave your comments and suggestions to help us improve it.
If you encounter any issue, we suggest you first search the issues section to see whether a solution already exists, or open up a new one otherwise.

We’d like to acknowledge the contributions made by our [developers and contributors](https://github.com/pearlnote/pearlnote/graphs/contributors) to
this project. Pearlnote won’t exist without your hard work. Your help is much appreciated.

## 7. Join us

Please feel free to fork this repository and contribute back using [pull requests](https://github.com/pearlnote/pearlnote/pulls).

If you find any problems or have any good ideas, feature requests, please submit here [issues](https://github.com/pearlnote/pearlnote/issues).


## 8. Donation

If you like our product, consider supporting us via [donate us](http://pearlnote.org/#donate).
We acknowledge the donations made by all the [donators](http://pearlnote.pearlnote.com/post/pearlnote-donation-list).

## 9. Related projects

* [Pearlnote Desktop App](https://github.com/pearlnote/desktop-app), [Download](http://app.pearlnote.com)
* [Pearlnote iOS](https://github.com/pearlnote/pearlnote-ios), [Download From App Store](https://itunes.apple.com/app/pearlnote/id1022302858)
* [Pearlnote Android](https://github.com/pearlnote/pearlnote-android), development phase

You are welcome to join us.

## 10. Contacts

* Email: pearlnote@pearlnote.com
* [Pearlnote BBS](http://bbs.pearlnote.com)
* [Pearlnote Google Group](https://groups.google.com/forum/#!forum/pearlnote)
* QQ Groups: 326073529, 256076853, 158716820

-----------------------------------------------------------------------

# 珠玑笔记

## 1. 介绍

珠玑笔记, 不只是笔记!

**特性**

* 高效笔记：珠玑笔记 有易操作的界面, 包含一款富文本编辑器和Markdown编辑器，让您的笔记记录更轻松和高效。对高阶用户，我们还提供`Vim` 和`Emacs` 编辑模式，助推你的写作速度更上层楼。
* 知识管理:  珠玑笔记 灵活而强大的“笔记本-笔记-标签”系统，让它成为你个人知识管理的利器。
* 分享: 你可以通过珠玑笔记同好友分享知识、想法和经历, 邀请好友加入你的笔记簿，通过云端交流信息。
* 协作: 珠玑笔记协助你与同事之间相互协作，激荡新思路，随时随地头脑风暴。
* 博客: 珠玑笔记也可以作为你的个人博客, 把你的知识传播的更远!

**其它特性**

* 支持Markdown编辑
* 写作模式
* `Vim` 及 `Emacs` 编辑模式
* 支持PDF导出
* 支持批量操作
* 博客自定义主题, 实现高度定制化

## 2. 为什么我们要创建珠玑笔记?
我们都曾是`Evernote`的忠实粉丝, 一直以来`Evernote`都是我们日常知识管理的有效工具。于是我们决定重新创造一款工具，提供`Evernote`所能提供的功能，同时弥补`Evernote`的不足，比如：
* 功能更强的文本编辑器：`Evernote`的编辑器不能满足我们的需求, 不能实现文档导航、不能贴代码(格式会乱掉, 作为程序员, 代码是我们的基本需求啊), 图片不能缩放等。
* `Evernote` 不支持所有人都喜爱的markdown语法，于是我们为珠玑笔记配备了一款可以实时渲染的markdown编辑器。
* 如果你是一名开发者，觉得手指怀念`Vim` 或 `Emacs` 了，那么我们还提供给你`Vim` 和 `Emacs` 写作模式，配合*markdown*的格式编辑，写作的时候再也不用去碰鼠标了。
* 知识积累和知识分享同样重要，因此大家都有自己的笔记账号和社交账号。但为什么这两者不能合二为一呢? 珠玑笔记 做到了将二者无缝衔接。
* 一套完整的、全平台覆盖的软件套装，包括了web、桌面、安卓、IOS设备，而且全部开源！
* 还有...

## 3. 获取珠玑笔记

珠玑笔记云笔记产品包括: 珠玑笔记 Web & Server(即本仓库), 桌面客户端, IOS, android. 4端全部开源!

如果想试用我们的产品，欢迎在 https://pearlnote.com 上注册, 珠玑笔记团队为你提供稳定可靠的服务。
担心服务厂商如何处理你的个人数据吗？你可以下载珠玑笔记安装在自己的服务器上, 通过珠玑笔记客户端连接与自建服务同步数据。

这里详细整理了珠玑笔记二进版和珠玑笔记开发版的安装教程, 请移步至:

* 珠玑笔记二进制详细安装教程:
    * [Windows](https://github.com/pearlnote/pearlnote/wiki/Pearlnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B---Windows)
    * [Mac, Linux](https://github.com/pearlnote/pearlnote/wiki/pearlnote%E4%BA%8C%E8%BF%9B%E5%88%B6%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)
* 珠玑笔记源码详细安装教程:
    <!-- * [Windows](https://github.com/pearlnote/pearlnote/wiki/Pearlnote-%E6%BA%90%E7%A0%81%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B----Windows) -->
    * [Mac, Linux](https://github.com/pearlnote/pearlnote/wiki/pearlnote%E5%BC%80%E5%8F%91%E7%89%88%E8%AF%A6%E7%BB%86%E5%AE%89%E8%A3%85%E6%95%99%E7%A8%8B)

## 4. 相关文档

更多详细的安装说明、问题处理和配置说明文档，请查看 [wiki](https://github.com/pearlnote/pearlnote/wiki)。

## 5. 如何对珠玑笔记进行二次开发

如果您有兴趣基于珠玑笔记二次开发，请查看 [How-to-develop-珠玑笔记](https://github.com/pearlnote/pearlnote/wiki/How-to-develop-pearlnote-%E5%A6%82%E4%BD%95%E5%BC%80%E5%8F%91pearlnote)。

## 6. 贡献者

在此对向珠玑笔记贡献力量的[贡献者们](https://github.com/pearlnote/pearlnote/graphs/contributors) 表示感谢。珠玑笔记因有你们而更完美!

## 7. 加入我们

欢迎提交[pull requests](https://github.com/pearlnote/pearlnote/pulls) 到珠玑笔记。

有任何问题或建议, 请先搜索[issue](https://github.com/pearlnote/pearlnote/issues)区是否已经有解决方法。如果没有，欢迎提交新issue。

珠玑笔记还有很多问题, 如果你喜欢它, 欢迎加入我们一起完善珠玑笔记。

## 8. 捐赠

如果您喜欢我们的产品，请考虑支持我们, [捐赠珠玑笔记](http://pearlnote.org/#donate)。

感谢[这些捐赠者](http://pearlnote.pearlnote.com/post/pearlnote-donation-list), 谢谢你们的鼓励, 珠玑笔记会一直坚持!

## 9. 其它相关项目
* [珠玑笔记 Desktop App](https://github.com/pearlnote/desktop-app), [下载地址](http://app.pearlnote.com)
* [珠玑笔记 iOS](https://github.com/pearlnote/pearlnote-ios), [从App Store下载](https://itunes.apple.com/zn/app/pearlnote/id1022302858?mt=8)
* [珠玑笔记 Android](https://github.com/pearlnote/pearlnote-android), 开发阶段

**数据库工具**:
* [Database Abstraction Layer](app/db/) - 统一的MongoDB和PostgreSQL接口
* [Migration Tools](tools/migration/) - 双向数据迁移工具

## 11. Docker 部署

珠玑笔记 支持 Docker 部署，并提供 MongoDB 和 PostgreSQL 两种数据库配置。

### 镜像构建

从源码构建：

```bash
docker build -t pearlnote:v2.7.0 .
```

或者直接使用预构建的镜像：

```bash
docker pull raptor/pearlnote:v2.7.0
```

### PostgreSQL 部署（推荐）

PostgreSQL 配置利用了 珠玑笔记 的数据库抽象层，推荐新部署使用：

```bash
# 启动 PostgreSQL 版本
docker-compose -f docker-compose.postgres.yml up -d

# 初始化数据库（首次运行）
docker exec postgres psql -U pearlnote -d pearlnote -f /docker-entrypoint-initdb.d/init.sql

# 或者使用迁移工具从 MongoDB 导入数据
docker exec pearlnote /opt/pearlnote/scripts/migrate_from_mongodb.sh
```

### MongoDB 部署

传统的 MongoDB 部署方式：

```bash
# 初始化 MongoDB 数据（首次运行）
docker run -d -v ~/mongo_data:/data --name=mongodb mongo:4.2
docker run -it --rm -v ~/pearlnote/mongodb_backup/pearlnote_install_data:/root/initdata --link=mongodb --entrypoint="" mongo:4.2 mongorestore -h mongodb -d pearlnote --dir /root/initdata/
docker rm -f mongodb

# 启动 MongoDB 版本
docker-compose -f docker-compose.mongodb.yml up -d
```

### 配置说明

配置文件 `conf/app.conf` 包含 PostgreSQL 和 MongoDB 两种数据库配置：

- 默认使用 PostgreSQL（`db.type=postgresql`）
- 如需使用 MongoDB，修改 `db.type=mongodb` 并取消注释 MongoDB 配置部分

自定义配置：
1. 修改 `conf/app.conf` 中的数据库类型和连接参数
2. 修改 `app.secret` 为你自己的值（安全必需）
3. PostgreSQL 部署使用 `docker-compose -f docker-compose.postgres.yml`
4. MongoDB 部署使用 `docker-compose -f docker-compose.mongodb.yml`

### 访问 珠玑笔记

部署完成后，访问 `http://localhost:9000` 即可使用 珠玑笔记。

欢迎加入我们!

## 联系&加入我们
* Email: pearlnote@pearlnote.com
* [珠玑笔记 社区](http://bbs.pearlnote.com)
* [QQ群](http://pearlnote.pearlnote.com/post/Pearlnote-groups)
* [珠玑笔记 Google Group](https://groups.google.com/forum/#!forum/pearlnote)
