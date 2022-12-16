<template>
  <div>
    <ManageHeader></ManageHeader>
    <div class="container">
      <nav class="columns">
        <div class="column">
          <div class="field is-grouped">
            <p class="control is-expanded">
              <input class="input is-small" type="text" v-model="SearchDatekey" placeholder="输入日期" @input="SearchDate">
            </p>
            <p class="control">
              <button class="button is-small is-info" :disabled="SearchDatekey.length > 0 ? false : true" @click="Clean">
                清空结果
              </button>
            </p>
          </div>
        </div>
      </nav>
      <nav class="columns">
        <div class="column">
          <div class="field is-grouped">
            <p class="control is-expanded">
              <input class="input is-small" type="text" v-model="SearchKey" placeholder="输入帐号" @input="Search">
            </p>
            <p class="control">
              <button class="button is-small is-info" :disabled="SearchKey.length > 0 ? false : true" @click="Clean">
                清空结果
              </button>
            </p>
          </div>
        </div>
      </nav>
      <div class="card events-card" v-if="loading">
        <header class="card-header">
          <p class="card-header-title">
            努力加载中.....
          </p>
        </header>
        <div class="content has-text-centered	min-heights" style="min-height: 11.3rem">
            <div class="com__box" v-if="loading" :style="loading? 'margin-top:5rem':''">
              <LoadIng></LoadIng>
            </div>
        </div>
      </div>
      <div v-else>
        <div class="card events-card" v-if="data.length <= 0">
          <header class="card-header">
            <p class="card-header-title">
              空空如也
            </p>
          </header>
          <div class="card-content">
            <div class="content has-text-centered	min-heights" style="min-height: 11.3rem">
              <EmptyEd></EmptyEd>
            </div>
          </div>
        </div>
        <div class="card events-card mt-5"  v-else>
          <header class="card-header">
            <p class="card-header-title">
              <span v-if="SearchKey.length == 0">{{title}}卡号列表-日期：<span class="has-text-link-dark mr-5">{{datelist[currentPage]}}</span></span> 帐号总数：<span class="has-text-danger-dark">{{data.length}}</span>
            </p>
              <div class="buttons are-small card-header-icon">
                <button class="button is-success" @click="copyIt">
                  复制本页数据到剪切板
                </button>
              </div>
          </header>
          <div class="card-content">
            <div class="content has-text-centered	min-heights" style="min-height: 11.3rem">
              <table class="table is-striped is-hoverable is-fullwidth is-narrow has-text-left">
                <thead>
                  <tr>
                    <td>编号</td>
                    <td>帐号</td>
                    <td>密码</td>
                    <td>
                      <span style="cursor: pointer" @click="sortTable">
                        今日金币
                        <span class="icon is-small">
                          <i class="fa" :class="goldSort?'fa-angle-up':'fa-angle-down'"></i>
                        </span>
                      </span>
                    </td>
                    <td>炮台倍数</td>
                    <td>金币截图</td>
                    <td>提号人</td>
                    <td>
                      <span style="cursor: pointer" @click="sortDatetime">
                        日期
                        <span class="icon is-small">
                          <i class="fa" :class="dateSort?'fa-angle-up':'fa-angle-down'"></i>
                        </span>
                      </span>
                    </td>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(it,index) in data" :key="it.id" @mouseover="()=>setBackage(true,it.cover)" @mouseout="()=>setBackage(false,'')">
                    <td>{{index + 1}}</td>
                    <td>{{it.account}}</td>
                    <td>{{it.password}}</td>
                    <td>{{makeNumber(it.today)}}</td>
                    <td>{{it.multiple}}</td>
                    <td><DefaultImg :img-url="rootUrl + it.cover" img-style="thisimg" ></DefaultImg></td>
                    <td>{{it.username}}</td>
                    <td><FormaTime :DateTime="it.updatetime"></FormaTime></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <PaginAtion v-if="data.length > 0 && SearchKey.length == 0" :total="total" :number="pageNumber" :GetData="makePageData"></PaginAtion>
      <NotIfication
        :showData="openerr" />
    </div>
    <div class="is-img" :style="{top:hoverTop+'px'}" v-if="imghover"><DefaultImg :img-url="rootUrl + currentImg"></DefaultImg></div>
  </div>
</template>
<script>
import { reactive, toRefs, defineComponent, onBeforeMount } from 'vue'
import { useRouter } from 'vue-router'
import useClipboard from 'vue-clipboard3'
const { toClipboard } = useClipboard()

import ManageHeader from '@/components/Other/Header'
import LoadIng from '@/components/Other/Loading'
import EmptyEd from '@/components/Other/Empty'
import PaginAtion from '@/components/Other/PaginAtion'
import FormaTime from '@/components/Other/FormaTime'
import NotIfication from "@/components/Other/Notification"
import DefaultImg from '@/components/Other/DefaultImg'


import Fetch from '@/helper/fetch'
import Config from '@/helper/config'
export default defineComponent({
  name: 'AccdList',
  components: { ManageHeader, LoadIng, EmptyEd, PaginAtion, FormaTime, NotIfication, DefaultImg },
  setup() {
    const router = useRouter()
    let states = reactive({
      loading: false,
      temp: [],
      data: [],
      datelist: [],
      SearchTemp:[],
      total: 0,
      rootUrl: Config.RootU,
      page: [],
      title: "",
      SearchKey: "",
      currentImg: '',
      imghover: false,
      hoverTop: 0,
      pageNumber: 1,
      currentPage: 0,
      TempCurrentPage: 0,
      SearchDatekey: "",
      goldSort: false,
      dateSort: false,
      openerr: {
        active: false,
        message: "",
        color: "",
        newtime: 0,
      }
    })
    onBeforeMount(async()=>{
      states.loading = true
      states.title = "已取"
      document.title = `${Config.GlobalTitle}-${states.title}帐号列表`
      const token = localStorage.getItem("token")
      const data = await Fetch(Config.Api.checklogin, {}, "get", token)
      if (data.status == 0) {
        GetData()
      }else {
        localStorage.removeItem("token")
        router.push("/")
      }
    })
    const GetData = async() => {
      const page = states.currentPage
      let url = Config.Api.datetime
      let date = states.datelist[page]
      if ( date == undefined ) date = ""
      url = `${url}?status=1&date=${date}`
      const token = localStorage.getItem("token")
      states.loading = true
      const d = await Fetch(url, {}, 'GET', token)
      if (d.status == 1) {
        states.data = d.data
        states.datelist = d.datelist
        // const NewData = makData(d.data)
        // states.SearchTemp = d.data
        states.temp = d.data
        states.total = d.datelist.length
        states.currentPage = page
        states.loading = false
      }else{
        states.data = []
        states.total = 0
        states.page = []
        states.loading = false
      }
    }
    const Search = async() => {
      const search = states.SearchKey
      if (search) {
        let url = Config.Api.search
        let date = search
        if ( date == undefined ) date = ""
        url = `${url}?status=1&key=${date}`
        const token = localStorage.getItem("token")
        states.loading = true
        const d = await Fetch(url, {}, 'GET', token)
        if (d.status == 1) {
          states.data = d.data
          states.loading = false
        }else{
          states.data = []
          states.total = 0
          states.page = []
          states.loading = false
        }
      }else {
        states.data = states.temp 
      }
    }
    const Clean = () => {
      states.SearchKey = ""
      states.SearchDatekey = ""
      states.data = states.temp
      states.currentPage = states.TempCurrentPage
    }

    const makeNumber = (n) =>{
      let x = "0"
      if ((n+"").length >= 9 && n !== 0) {
        let a = n / 100000000
        a = a.toFixed(1)
        x = `${a}亿`
      }else if (n === 123) {
        x = "识别错误"
      }else{
        if (n !== 0 ) {
          let a = n / 10000
          a = a.toFixed(1)
          x = `${a}万`
        }
      }
      return x
    }
/*
    const FormatTime = (timestamp) => {
      const date = new Date(timestamp * 1000);
      const Y = date.getFullYear()
      const M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1)
      const D = date.getDate()
      return `${Y}-${M}-${D}`
    }

    const makData = (data) => {
      let date = []
      let newData= []
      data.forEach(el => {
        date = [...date, FormatTime(el.updatetime)]
      })
      date = [... new Set(date)]
      date = quickSort(date)
      date.forEach((e) => {
        let newDate = {
          date: e,
          data: []
        }
        data.forEach(el => {
          if (e == FormatTime(el.updatetime)) {
            newDate.data = [...newDate.data, el]
          }
        })
        newData = [...newData, newDate]
      })
      return newData
    }

    const quickSort = (arr, s = true) => {
      if (arr.length <= 1) {
        return arr
      }
      let pivotIndex = Math.floor(arr.length / 2)
      let pivot = arr.splice(pivotIndex, 1)[0]
      let left = []
      let right = []
      for (let i = 0; i < arr.length; i++) {
        if (s) {
          if (arr[i] > pivot) {
            left.push(arr[i])
          } else {
            right.push(arr[i])
          }
        }else {
          if (arr[i] < pivot) {
            left.push(arr[i])
          } else {
            right.push(arr[i])
          }
        }
      }
      return quickSort(left).concat([pivot], quickSort(right))
    }
*/
    const setBackage = (hover,img) => {
      states.currentImg = img
      states.imghover = hover
      let  scrollTop = document.documentElement.scrollTop || window.pageYOffset || document.body.scrollTop;
      states.hoverTop = scrollTop + 60
    }
    const makePageData = (page) => {
      states.currentPage = page - 1
      states.SearchKey = ""
      states.SearchDatekey = ""
      GetData(page)
    }
    const SearchDate = async() => {
      const search = states.SearchDatekey
      if (search) {
        let url = Config.Api.datetime
        let date = search
        if ( date == undefined ) date = ""
        url = `${url}?status=1&date=${date}`
        const token = localStorage.getItem("token")
        states.loading = true
        const d = await Fetch(url, {}, 'GET', token)
        let page = 0
        states.datelist.map((el,i) => {
          if (el == search) page = i
        })
        if (d.status == 1) {
          states.data = d.data
          states.datelist = d.datelist
          // console.log(page)

          states.TempCurrentPage = states.currentPage
          // const NewData = makData(d.data)
          // states.SearchTemp = d.data
          states.total = d.datelist.length
          states.currentPage = page
          states.loading = false
        }else{
          states.data = []
          states.total = 0
          states.page = []
          states.loading = false
        }
      }
    }
    const sortTable = () => {
      states.goldSort = !states.goldSort
      states.data.sort((a, b)=>{
          if (states.goldSort) {
            return a.today - b.today
          }else{
            return b.today - a.today
          }
        })
    }
    const sortDatetime = () => {
      states.dateSort = !states.dateSort
      states.data.sort((a, b)=>{
          if (states.dateSort) {
            return a.updatetime - b.updatetime
          }else{
            return b.updatetime - a.updatetime
          }
        })
    }
    const makeNumberINT = (n) =>{
      let x = "0"
      if ((n+"").length >= 9 && n !== 0) {
        const a = Math.floor(n / 100000000)
        x = `${a}`
      }else if (n === 123) {
        x = "识别错误"
      }else{
        if (n !== 0 ) {
          const a = Math.floor(n / 10000)
          x = `${a}`
        }
      }
      return x
    }
    const ShowMessage = (e) => {
      states.openerr = e
    }
    const copyIt = async() => {
      const d = states.data
      let t = []
      d.map(el=>{
        t = [...t, `${el.account}\t${el.password}\t${makeNumberINT(el.today)}\t${el.multiple}`]
      })
      const dd = t.join("\r\n")
      await toClipboard(dd)
      const e = {
        active: true,
        message: "复制帐号成功",
        color: "is-success"
      }
      ShowMessage(e)
    }

    return {
      ...toRefs(states),
      GetData,
      Search,
      Clean,
      makeNumber,
      setBackage,
      makePageData,
      SearchDate,
      sortTable,
      sortDatetime,
      copyIt
    }
  },
})
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
