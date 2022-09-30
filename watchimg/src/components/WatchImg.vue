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
      <div class="container" v-if="list !== null">
        <div class="content" v-for="(key,index) in list" :key="key">
          <h4>{{index}}</h4>
          <div class="container is-mobile">
            <div class="table-container">
              <table class="table is-hoverable is-fullwidth">
                <tbody>
                  <tr v-for="(ix,imx) in key" :key="imx">
                    <td @mouseover="setBackage(true,ix.imgurl)" @mouseout="setBackage(false,'')">{{ix.title}}</td>
                    <td>{{ix.gold}}亿</td>
                    <td><img class="thisimg" :src="rootUrl + ix.imgurl" /></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <hr />
        </div>
      </div>
      <div class="is-img" :style="{top:hoverTop+'px'}" v-if="hover"><img :src="rootUrl + currentImg" /></div>
    </div>
  </section>
</template>

<script>
  const rootUrl = "http://127.0.0.1:13002"
  export default {
    name: 'WatchImg',
    data(){
      return{
        list: {
          data:[]
        },
        rootUrl: rootUrl,
        currentImg: '',
        hover: false,
        hidden: false,
        code: '',
        error: false,
        vgh8MOC: "avcTd$auZFNJ"
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
      async GetData(){
        const url = `${rootUrl}/api/data`
        let data = await this.getData(url)
        
        // console.log(data)
        let a = this.json_sort(data)
        this.list = a
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
      json_sort(data){
        data = data || '{}';
        let arr = [];
        for (var key in data) {
          arr.push(key)
        }
        arr.reverse();
        let str = '{';
        for (var i in arr) {
          const d = JSON.stringify(data[arr[i]])
          str += '"' + arr[i] + '":' + d + ','
        }
        str = str.substr(0, str.length - 1)
        str += '}'
        str = JSON.parse(str);
        return str;
      },
      getData(url){
        let requestConfig = {
          method: 'get',
          headers: {
            Accept: 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
            'Content-type':'text/html;charset=UTF-8'
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