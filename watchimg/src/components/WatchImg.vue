<template>
    <section class="section">
      <div class="container" v-if="!hidden">
        <section class="container w350">
          <div class="field">
            <p class="control">
              <input class="input" type="username" placeholder="用户名" v-model="username">
            </p>
          </div>
          <div class="field">
            <p class="control">
              <input class="input" type="password" placeholder="密码" v-model="password">
            </p>
            <p class="help is-danger" v-if="error">登陆错误</p>
          </div>
          <div class="field">
            <p class="control">
              <button class="button is-success" @click="check">
                登陆
              </button>
            </p>
          </div>
        </section>
      </div>
      <div v-else>
      <div class="container" v-if="list !== null && list.length > 0 && status === 1">
        <div class="content">
          <div class="container is-mobile">
            <div class="table-container">
              <table class="table is-hoverable is-fullwidth">
                <thead>
                  <tr>
                    <td><label class="checkbox"><input type="checkbox" @click="checkall" />全选</label></td>
                    <td colspan="5"></td>
                    <td>
                      <div class="buttons are-small">
                        <button class="button is-info" @click="getall" :disabled="temp.length <= 0">
                          提取选中帐号
                        </button>
                      </div>
                    </td>
                  </tr>
                  <tr>
                    <td>选择</td>
                    <td>帐号</td>
                    <td>今日金币</td>
                    <td>昨日金币</td>
                    <td>金币截图</td>
                    <td>日期</td>
                    <td>操作</td>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(item) in list" :key="item.id" @mouseover="setBackage(true,item.cover)" @mouseout="setBackage(false,'')">
                    <td><input type="checkbox" v-model="item.check" @click="(e)=>checkBox(e,item.account)"></td>
                    <td>{{item.account}}</td>
                    <td>{{makeNumber(item.today)}}</td>
                    <td>{{makeNumber(item.yesterday)}}</td>
                    <td><img class="thisimg" :src="rootUrl + item.cover" /></td>
                    <td>{{FormatTime(item.datetime)}}</td>
                    <td>
                      <div class="buttons are-small">
                        <button class="button is-info" @click="delit(item.account)">
                          提取帐号
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <hr />
          <nav class="pagination buttons" role="navigation" aria-label="pagination">
            <button class="button pagination-previous" @click="setPreviousPage" title="This is the first page" :disabled="current !== 1 && page.length > 1?false:true">上一页</button>
            <button class="button pagination-next" @click="setNextPage" :disabled="current !== page.length && page.length > 1?false:true">下一页</button>
            <ul class="pagination-list">
              <li v-for="(item) in page" :key="item">
                <a class="pagination-link" @click="setPage(item)" :class="item === current ? 'is-current': ''" :aria-label="'Page '+item" :aria-current="item === current ? 'page': ''">{{item}}</a>
              </li>
            </ul>
          </nav>
        </div>
      </div>
    </div>

    <div class="modal is-active" v-if="loading"><div class="modal-background"></div></div>

    <div class="modal is-active" v-if="hasaccount">
      <div class="modal-background"></div>
      <div class="modal-card">
        <header class="modal-card-head">
          <p class="modal-card-title">帐号信息</p>
          <button class="delete" aria-label="close" @click="closeMo"></button>
        </header>
        <section class="modal-card-body">
          <p>帐号：{{account.accountgs}}</p>
          <p>密码：{{account.password}}</p>
        </section>
        <footer class="modal-card-foot">
          <button class="button is-success" @click="copyAccount">复制到剪切板</button>
          <button class="button" @click="closeMo">关闭</button>
        </footer>
      </div>
    </div>
    <div class="modal is-active" v-if="haslist">
      <div class="modal-background"></div>
      <div class="modal-card">
        <header class="modal-card-head">
          <p class="modal-card-title">帐号信息</p>
          <button class="delete" aria-label="close" @click="closeMo"></button>
        </header>
        <section class="modal-card-body">
          <p class="autohidden">{{listcount}}</p>
        </section>
        <footer class="modal-card-foot">
          <button class="button is-success" @click="copyList">复制到剪切板</button>
          <button class="button" @click="closeMo">关闭</button>
        </footer>
      </div>
    </div>
    <div class="notification is-danger error" :style="{top:hoverTop+'px'}" v-if="hover">
      <button class="delete" @click="closErr"></button>
      <p>删除帐号出错！</p>
    </div>
    <div class="notification is-success error" :style="{top:hoverTop+'px'}" v-if="ok">
      <button class="delete" @click="closErr"></button>
      <p>复制帐号成功</p>
    </div>
    <div class="is-img" :style="{top:hoverTop+'px'}" v-if="imghover"><img :src="rootUrl + currentImg" /></div>
  </section>
</template>

<script>
  import useClipboard from 'vue-clipboard3'

  const { toClipboard } = useClipboard()
  const rootUrl = "http://127.0.0.1:13002"
  export default {
    name: 'WatchImg',
    data(){
      return{
        list: [],
        rootUrl: rootUrl,
        currentImg: '',
        hover: false,
        imghover: false,
        hidden: false,
        username: '',
        error: false,
        password: "",
        total: 0,
        current: 1,
        status: 0,
        page:[],
        hasaccount: false,
        account: {},
        atest: "",
        ok: false,
        temp:[],
        checkbos:[],
        hoverTop: 0,
        listcount: "",
        haslist: false,
        loading: false,
      }
    },
    created(){
      let key = localStorage.getItem("token")
      if(key !== null) {
        this.GetData()
        this.hidden = true
        this.error = false
      }
    },
    methods: {
      cleanAll(){
        this.list = []
        this.total = 0
        this.status = 0
        this.page = []
        localStorage.removeItem("token")
        this.hidden = false
        this.temp = []
      },
      setBackage(hover,img){
        this.currentImg = img
        this.imghover = hover
        let  scrollTop = document.documentElement.scrollTop || window.pageYOffset || document.body.scrollTop;
        this.hoverTop = scrollTop
      },
      async getall(){
        const list = this.temp
        if (list.length > 0 ) {
          const url = `${rootUrl}/api/delist`
          this.loading = true
          const data = await this.postListData(url)
          // console.log(data)
          if (data.status == -1) {
            this.cleanAll()
          }else if (data.status == 1) {
            this.openError()
          }else {
            if (data.list !== null && data.list.length > 0) {
              // data.list.forEach((y) => {
              //   const z = JSON.parse(y)
              //   this.listcount += `${z.accountgs}----${z.password}----${z.other}\r\n`
              // })
              let xyz = []
              data.list.forEach((el) => {
                let what = ""
                const z = JSON.parse(el)
                this.list.forEach((y) => {
                  if (el.accountgs !== y.account) {
                    what = `${z.accountgs}----${z.password}----${this.makeNumber(y.today)}`
                  }
                })
                xyz = [...xyz, what]
              })
              this.list = this.list.filter((el) => {
                return data.list.every((el2) => {
                  const eel = JSON.parse(el2)
                  return eel.accountgs !== el.account
                })
              })
              // let what = data.list.filter(function(item1) {
              //   let xyz = JSON.parse(item1)
              //   return this.list.every(function(item2) {
              //     return item2.account !== xyz.accountgs
              //   })
              // })
              if (this.list.length === 0) {
                location.reload()
              }else {
                this.listcount = xyz.join("\r\n")
                this.haslist = true
                this.loading = false
                this.temp = []
              }
            }else{
              this.loading = false
              this.openError()
            }
          }
        }else{
          this.loading = false
          this.openError()
        }
      },
      async delit(account){
        const url = `${rootUrl}/api/delone`
        const data = await this.delData(url,account)
        // console.log(data)
        if (data.status == -1) {
          this.list = []
          this.total = 0
          this.status = 0
          this.page = []
          localStorage.removeItem("token")
          this.hidden = false
        }else if (data.status == 1) {
          this.openError()
        }else {
          if (data.data !== "{}") {
            const x = JSON.parse(data.data)
            let d = ""
            this.list = this.list.filter((e) => {
              d = `${x.accountgs}----${x.password}----${this.makeNumber(e.today)}`
              return e.account !== account
            })
            this.account = x
            this.atest = d
            this.hasaccount = true
            this.temp = []
          }else{
            this.openError()
          }
        }
      },
      async copyAccount(){
        const _this = this
        await toClipboard(_this.atest)
        this.onCopy()
      },
      async copyList() {
        const _this = this
        await toClipboard(_this.listcount)
        this.onCopy()
      },
      openError(){
        let  scrollTop = document.documentElement.scrollTop || window.pageYOffset || document.body.scrollTop;
        this.hoverTop = scrollTop
        this.hover = true
        const _this = this
        setTimeout(function(){
          _this.hover = false
        },10000)
      },
      closErr(){
        this.hover = false
      },
      closeMo(){
        this.hasaccount = false
        this.haslist = false
      },
      onCopy(){
        let  scrollTop = document.documentElement.scrollTop || window.pageYOffset || document.body.scrollTop;
        this.hoverTop = scrollTop
        this.ok = true
        const _this = this
        setTimeout(function(){
          _this.ok = false
        },10000)
      },
      setPage(p){
        if (p !== this.current && p >= 1) {
          this.current = p
          this.GetData()
          this.temp = []
          // console.log(this.current)
        }
      },
      setNextPage() {
        if (this.current !== this.page.length) {
          this.temp = []
          this.current = this.current + 1
          this.GetData()
          // console.log(this.current)
        }
      },
      setPreviousPage(){
        if (this.current > 1) {
          this.temp = []
          this.current = this.current - 1
          this.GetData()
          // console.log(this.current)
        }
      },
      makeNumber(n){
        let x = "0"
        if ((n+"").length >= 9 && n !== 0) {
          const a = n / 100000000
          x = `${a}亿`
        }else if (n === 123) {
          x = "识别错误"
        }else{
          if (n !== 0 ) {
            const a = n / 10000
            x = `${a}万`
          }
        }
        return x
      },
      FormatTime(timestamp) {
        var date = new Date(timestamp * 1000);
        var Y = date.getFullYear() + '-'
        var M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1) + '-'
        var D = date.getDate() + ' '
        return Y+M+D
      },
      async GetData(){
        const url = `${rootUrl}/api/data`
        let data = await this.getData(url)
        if (data.status === 1) {
          this.list = data.data.map((e)=>{
            e.check = false
            return e
          })
          // console.log(this.list)
          this.total = data.total
          this.status = data.status
          this.page = this.makePage(data.total)
          // console.log(this.page)
        }else{
          this.list = []
          this.total = 0
          this.status = 0
          this.page = []
          localStorage.removeItem("token")
          this.hidden = false
        }
      },
      async check(){
        let key = localStorage.getItem("token")
        if (key === null) {
          if (this.username.length > 0 && this.password.length > 0) {
            const url = `${rootUrl}/api/login`
            let data = await this.postData(url)
            if (data.status == 1) {
              this.error = true
            }else {
              this.token = data.token
              localStorage.setItem("token", data.token)
              this.GetData()
              this.hidden = true
              this.error = false
            }
          }else{
            this.error = true
          }
          // localStorage.setItem("code", this.vgh8MOC)
        }else{
          this.GetData()
          this.hidden = true
          this.error = false
        }
      },
      makePage(t){
        let x = []
        const p = Math.ceil(t/100)
        for (let i = 0; i < p; i++) {
          x = [...x, i + 1]
        }
        return x
      },
      getData(url){
        let key = localStorage.getItem("token")
        let requestConfig = {
          method: 'get',
          headers: {
            Accept: 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${key}`,
          }
        }
        const page = this.current
        const newUrl = `${url}?page=${page}`
        return new Promise((resolve) => {
          fetch(newUrl, requestConfig)
            .then(res => {
              if(res.ok) {
                return res.blob()
              }else {
                resolve({
                  status: 0,
                  message: res.status
                })
              }
            })
            .then(blob => {
              var reader = new FileReader();
              reader.onload = function () {
                var text = reader.result;
                // console.log(text)
                // const json = makeData(pages, text)
                resolve(JSON.parse(text))
              }
              reader.readAsText(blob, 'UTF-8')
            })
            .catch((err) => {
              resolve({
                status: 1,
                message: err.message
              })
            })
        })
      },
      postData(url){
        const body = {
          'username': this.username,
          'password': this.password
        }
        
        let requestConfig = {
          method: 'POST',
          body: JSON.stringify(body),
          headers: {
            Accept: 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
            'Content-Type': 'application/json'
          }
        }
        return new Promise((resolve) => {
          fetch(url, requestConfig)
            .then(res => {
              if(res.ok) {
                return res.blob()
              }else {
                resolve({
                  status: 1,
                  message: res.status
                })
              }
            })
            .then(blob => {
              var reader = new FileReader();
              reader.onload = function () {
                var text = reader.result;
                // console.log(pages)
                // const json = makeData(pages, text)
                resolve(JSON.parse(text))
              }
              reader.readAsText(blob, 'UTF-8')
            })
            .catch((err) => {
              resolve({
                status: 1,
                message: err.message
              })
            })
        })
      },
      postListData(url){
        let key = localStorage.getItem("token")
        const body = {
          'list': this.temp
        }
        
        let requestConfig = {
          method: 'DELETE',
          body: JSON.stringify(body),
          headers: {
            Accept: 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${key}`,
          }
        }
        return new Promise((resolve) => {
          fetch(url, requestConfig)
            .then(res => {
              if(res.ok) {
                return res.blob()
              }else {
                resolve({
                  status: 1,
                  message: res.status
                })
              }
            })
            .then(blob => {
              var reader = new FileReader();
              reader.onload = function () {
                var text = reader.result;
                // console.log(pages)
                // const json = makeData(pages, text)
                resolve(JSON.parse(text))
              }
              reader.readAsText(blob, 'UTF-8')
            })
            .catch((err) => {
              resolve({
                status: 1,
                message: err.message
              })
            })
        })
      },
      delData(url,account){
        let key = localStorage.getItem("token")
        const body = {
          'account': account
        }
        
        let requestConfig = {
          method: 'DELETE',
          body: JSON.stringify(body),
          headers: {
            Accept: 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${key}`,
          }
        }
        return new Promise((resolve) => {
          fetch(url, requestConfig)
            .then(res => {
              if(res.ok) {
                return res.blob()
              }else {
                resolve({
                  status: 1,
                  message: res.status
                })
              }
            })
            .then(blob => {
              var reader = new FileReader();
              reader.onload = function () {
                var text = reader.result;
                // console.log(pages)
                // const json = makeData(pages, text)
                resolve(JSON.parse(text))
              }
              reader.readAsText(blob, 'UTF-8')
            })
            .catch((err) => {
              resolve({
                status: 1,
                message: err.message
              })
            })
        })
      },
      checkBox(e, account){
        if (e.target.checked) {
          this.temp = [...this.temp, account]
        }else{
          this.temp = this.temp.filter((el) => {
            if (el !== account) return el
          })
        }
      },
      checkall(e){
        if (e.target.checked) {
          this.temp = []
          this.list.forEach((el) => {
            this.temp =  [...this.temp,el.account]
          })
          this.list = this.list.map((el) => {
            el.check =  true
            return el
          })
        }else{
          this.temp = []
          this.list = this.list.map((el) => {
            el.check =  false
            return el
          })
        }
      }
    }
  }
</script>
<style scoped>
.thisimg {
  width: 150px;
  height: 30px;
  max-width: 150px;
  max-height: 30px;
  display: block;
}
.error {
  position:absolute;
  right: 0px;
  width: 240px;
  z-index: 9999999999;
}
.w350 {
  width: 350px;
  margin-top: 30px
}
.is-img {
  position: absolute;
  display: block;
  right: 10px;
  max-width: 340px;
  width: 340px;
}
.autohidden {
  max-height: 300px;
  overflow-x: auto;
}
</style>