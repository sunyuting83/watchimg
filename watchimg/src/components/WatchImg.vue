<template>
    <section class="section">
      <div class="container" v-if="!hidden">
        <div class="field has-addons">
          <div class="control">
            <input class="input" type="text" placeholder="验证码" v-model="code">
          </div>
          <div class="control">
            <a class="button is-info" @click="check">
              确定
            </a>
          </div>
          <p class="help is-danger" v-if="error">验证码错误</p>
        </div>
      </div>
      <div v-else>
      <div class="container" v-if="list !== null && list.length > 0 && status === 1">
        <div class="content">
          <div class="container is-mobile">
            <div class="table-container">
              <table class="table is-hoverable is-fullwidth">
                <thead>
                  <tr>
                    <td>帐号</td>
                    <td>当前金币</td>
                    <td>上次金币</td>
                    <td>金币截图</td>
                    <td>日期</td>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(item) in list" :key="item.id">
                    <td @mouseover="setBackage(true,item.cover)" @mouseout="setBackage(false,'')">{{item.account}}</td>
                    <td>{{makeNumber(item.today)}}</td>
                    <td>{{makeNumber(item.yesterday)}}</td>
                    <td><img class="thisimg" :src="rootUrl + item.cover" /></td>
                    <td>{{FormatTime(item.datetime)}}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <hr />
          <nav class="pagination" role="navigation" aria-label="pagination">
            <a class="pagination-previous" @click="setPreviousPage" title="This is the first page" :disabled="current !== 1 && page.length > 1?'true':'false'">上一页</a>
            <a class="pagination-next" @click="setNextPage" :disabled="current !== page.length && page.length > 1?'true':'false'">下一页</a>
            <ul class="pagination-list">
              <li v-for="(item) in page" :key="item">
                <a class="pagination-link" @click="setPage(item)" :class="item === current ? 'is-current': ''" :aria-label="'Page '+item" :aria-current="item === current ? 'page': ''">{{item}}</a>
              </li>
            </ul>
          </nav>
        </div>
      </div>
      <div class="is-img" :style="{top:hoverTop+'px'}" v-if="hover"><img :src="rootUrl + currentImg" /></div>
    </div>
  </section>
</template>

<script>
  const rootUrl = "http://192.168.1.90:13002"
  export default {
    name: 'WatchImg',
    data(){
      return{
        list: [],
        rootUrl: rootUrl,
        currentImg: '',
        hover: false,
        hidden: false,
        code: '',
        error: false,
        vgh8MOC: "avcTd$auZFNJ",
        total: 0,
        current: 1,
        status: 0,
        page:[]
      }
    },
    created(){
      let key = localStorage.getItem("code")
      if(key === this.vgh8MOC) {
        this.GetData()
        this.hidden = true
        this.error = false
      }
    },
    methods: {
      setBackage(hover,img){
        this.currentImg = img
        this.hover = hover
        let  scrollTop = document.documentElement.scrollTop || window.pageYOffset || document.body.scrollTop;
        this.hoverTop = scrollTop
      },
      setPage(p){
        if (p !== this.current && p >= 1) {
          this.current = p
          this.GetData()
          // console.log(this.current)
        }
      },
      setNextPage() {
        if (this.current !== this.page.length) {
          this.current = this.current + 1
          this.GetData()
          // console.log(this.current)
        }
      },
      setPreviousPage(){
        if (this.current > 1) {
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
        
        // console.log(data)
        if (data.status === 1) {
          this.list = data.data
          this.total = data.total
          this.status = data.status
          this.page = this.makePage(data.total)
          // console.log(this.page)
        }
      },
      check(){
        let key = localStorage.getItem("code")
        if (this.code === this.vgh8MOC) {
          this.GetData()
          this.hidden = true
          this.error = false
          if (key === null) localStorage.setItem("code", this.vgh8MOC)
        }else{
          this.error = true
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
        let requestConfig = {
          method: 'get',
          headers: {
            Accept: 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
            'Content-type':'text/html;charset=UTF-8'
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
      }
    }
  }
</script>
<style scoped>
.is-img {
  position: absolute;
  display: block;
  right: 10px;
  max-width: 752px;
  width: 752px;
}
.thisimg {
  width: 150px;
  height: 30px;
  max-width: 150px;
  max-height: 30px;
  display: block;
}
</style>