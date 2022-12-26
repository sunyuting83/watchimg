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
              <span v-if="SearchKey.length == 0">{{title}}卡号列表-日期：<span class="has-text-link-dark mr-5">{{datelist[currentPage]}}</span></span> 当日帐号总数：<span class="has-text-danger-dark">{{data.length}}</span>
            </p>
              <div class="buttons are-small card-header-icon">
                <button class="button is-warning is-light" @click="gotoAccount">金币排列</button>
              </div>
          </header>
          <div class="card-content">
            <div class="content has-text-centered	min-heights" style="min-height: 11.3rem">
              <table class="table is-striped is-hoverable is-fullwidth is-narrow has-text-left">
                <thead>
                  <tr>
                    <td><label class="checkbox"><input type="checkbox" @click="checkall" />全选</label></td>
                    <td colspan="8"></td>
                    <td>
                      <div class="buttons are-small">
                        <PopoButton
                          message="删除选中帐号" color="is-danger" :callBack="delAllSql"></PopoButton>
                        <button class="button is-info" @click="getall" :disabled="checkTemp.length <= 0">
                          提取选中帐号
                        </button>
                      </div>
                    </td>
                  </tr>
                  <tr>
                    <td>选择</td>
                    <td>序号</td>
                    <td>帐号</td>
                    <td>
                      <span style="cursor: pointer" @click="sortTable">
                        今日金币
                        <span class="icon is-small">
                          <i class="fa" :class="goldSort?'fa-angle-up':'fa-angle-down'"></i>
                        </span>
                      </span>
                    </td>
                    <td>昨日金币</td>
                    <td>炮台倍数</td>
                    <td>金币截图</td>
                    <td>
                      <span style="cursor: pointer" @click="sortExptime">
                        到期时间
                        <span class="icon is-small">
                          <i class="fa" :class="expdateSort?'fa-angle-up':'fa-angle-down'"></i>
                        </span>
                      </span>
                    </td>
                    <td>
                      <span style="cursor: pointer" @click="sortDatetime">
                        日期
                        <span class="icon is-small">
                          <i class="fa" :class="dateSort?'fa-angle-up':'fa-angle-down'"></i>
                        </span>
                      </span>
                    </td>
                    <td>操作</td>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(it,index) in data" :key="it.id" @mouseover="()=>setBackage(true,it.cover)" @mouseout="()=>setBackage(false,'')">
                    <td><input type="checkbox" v-model="it.check" @click="(e)=>checkBox(e,it.account)"></td>
                    <td>{{index + 1}}</td>
                    <td>{{it.account}}</td>
                    <td>{{makeNumber(it.today)}}</td>
                    <td>{{makeNumber(it.yesterday)}}</td>
                    <td>{{it.multiple}}</td>
                    <td><DefaultImg :img-url="rootUrl + it.cover" img-style="thisimg" ></DefaultImg></td>
                    <td><ExpTime :DateTime="it.expdate" /></td>
                    <td><FormaTime :DateTime="it.datetime" /></td>
                    <td>
                      <div class="buttons are-small">
                        <PopoButton
                          message="删除帐号" color="is-danger" :callBack="()=>{delToSql(item.account)}"></PopoButton>
                        <button class="button is-info" @click="delit(it.account)">
                          提取帐号
                        </button>
                      </div>
                    </td>
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
      <RenewalCard
        :showData="openModal"
        :ShowMessage="ShowMessage" />
      <ListData
        :showData="arrayModal"
        :Close="closeArray"
        :ShowMessage="ShowMessage" />
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
import ExpTime from '@/components/Other/ExpTime'
import RenewalCard from '@/components/Other/Renewal'
import ListData from '@/components/Other/ListData'
import NotIfication from "@/components/Other/Notification"
import PopoButton from '@/components/Other/PopoButton'
import DefaultImg from '@/components/Other/DefaultImg'


import Fetch from '@/helper/fetch'
import Config from '@/helper/config'
export default defineComponent({
  name: 'AccdDateList',
  components: { ManageHeader, LoadIng, EmptyEd, PaginAtion, FormaTime, NotIfication, RenewalCard, ListData, PopoButton, DefaultImg, ExpTime },
  setup() {
    const router = useRouter()
    let states = reactive({
      loading: false,
      temp: [],
      data: [],
      checkTemp: [],
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
      expdateSort: false,
      openerr: {
        active: false,
        message: "",
        color: "",
        newtime: 0,
      },
      openModal:{
        active: false,
        account: {
          password: "",
          accountgs: "",
        },
        atest: "",
      },
      arrayModal:{
        active: false,
        datacount: "",
      }
    })
    onBeforeMount(async()=>{
      states.loading = true
      states.title = "未提取"
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
      url = `${url}?status=0&date=${date}`
      const token = localStorage.getItem("token")
      states.loading = true
      const d = await Fetch(url, {}, 'GET', token)
      if (d.status == 1) {
        states.data = d.data.map((e)=>{
          e.check = false
          return e
        })
        states.datelist = d.datelist
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
        url = `${url}?status=0&key=${date}`
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

    const getnum = (num) =>  {
      const s = num.toString();
      const result = s.substring(0,s.indexOf(".")+2);
      return result
    }

    const makeNumber = (n) =>{
      let x = "0"
      if ((n+"").length >= 9 && n !== 0) {
        let a = n / 100000000
        x = `${getnum(a)}亿`
      }else if (n === 123) {
        x = "识别错误"
      }else{
        if (n !== 0 ) {
          let a = n / 10000
          x = `${getnum(a)}万`
        }
      }
      return x
    }
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
        url = `${url}?status=0&date=${date}`
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
          return a.datetime - b.datetime
        }else{
          return b.datetime - a.datetime
        }
      })
    }
    const sortExptime = () => {
      states.expdateSort = !states.expdateSort
      states.data.sort((a, b)=>{
        if (states.expdateSort) {
          return a.expdate - b.expdate
        }else{
          return b.expdate - a.expdate
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

    const gotoAccount = () => {
      router.push("acclist")
    }

    const checkBox = (e, account) => {
      if (e.target.checked) {
        states.checkTemp = [...states.checkTemp, account]
      }else{
        states.checkTemp = states.checkTemp.filter((el) => {
          if (el !== account) return el
        })
      }
    }
    const checkall = (e) => {
      if (e.target.checked) {
        states.checkTemp = []
        states.data.forEach((el) => {
          states.checkTemp =  [...states.checkTemp,el.account]
        })
        states.data = states.data.map((el) => {
          el.check =  true
          return el
        })
      }else{
        states.checkTemp = []
        states.data = states.data.map((el) => {
          el.check =  false
          return el
        })
      }
    }
    const delit = async(account) => {
      const token = localStorage.getItem("token")
      const params = {
        'account': account
      }
      const data = await Fetch(Config.Api.delone, params, "delete", token)
      const e = {
          active: true,
          message: data.message,
          color: "is-danger",
          newtime: 0,
        }
      if (data.status == -1) {
        states.data = []
        states.total = 0
        states.page = []
        localStorage.removeItem("token")
      }else if (data.status == 1) {
        states.data = states.data.filter((e)=> {
          return e.account != account
        })
        ShowMessage(e)
      }else {
        if (data.data !== "{}") {
          const x = data.data
          let d = ""
          states.data = states.data.filter((e) => {
            if (e.account == x.accountgs) {
              d = `${x.accountgs}\t${x.password}\t${makeNumberINT(e.today)}\t${e.multiple}`
            }
            return e.account !== x.accountgs
          })
          states.total = states.total - 1
          states.openModal.active = true
          states.openModal.account = x
          states.openModal.atest = d
          states.checkTemp = []
        }else{
          states.data = states.data.filter((e)=> {
            return e.account != account
          })
          ShowMessage(e)
        }
      }
    }
    const delToSql = async(account) => {
      const token = localStorage.getItem("token")
      const params = {
        'account': account
      }
      const data = await Fetch(Config.Api.delonesql, params, "delete", token)
      const e = {
          active: true,
          message: data.message,
          color: "is-danger",
          newtime: 0,
        }
      if (data.status == -1) {
        states.data = []
        states.total = 0
        states.page = []
        localStorage.removeItem("token")
      }else if (data.status == 1) {
        ShowMessage(e)
      }else {
        states.data = states.data.filter((e) => {
          return e.account !== account
        })
        states.total = states.total - 1
        const e = {
          active: true,
          message: data.message,
          color: "is-success",
          newtime: 0,
        }
        ShowMessage(e)
      }
    }

    const getall = async() => {
      const list = states.checkTemp
      const e = {
          active: true,
          message: "获取失败",
          color: "is-danger",
          newtime: 0,
        }
      if (list.length > 0 ) {
        states.loading = true
        const token = localStorage.getItem("token")
        const params = {
          'list': states.checkTemp
        }
        const data = await Fetch(Config.Api.delist, params, "DELETE", token)
        if (data.status == -1) {
          states.data = []
          states.total = 0
          states.page = []
          localStorage.removeItem("token")
        }else if (data.status == 1) {
          states.data = states.data.filter((el) => {
            return list.every((el2) => {
              return el2 !== el.account
            })
          })
          states.checkTemp = []
          states.loading = false
          ShowMessage(e)
        }else {
          if (data.list !== null && data.list.length > 0) {
            let xyz = []
            data.list.map((el) => {
              let what = ""
              const z = el
              states.data.map((y) => {
                if (z.accountgs == y.account) {
                  // console.log(z.accountgs,y.account)
                  what = `${z.accountgs}\t${z.password}\t${makeNumberINT(y.today)}\t${y.multiple}`
                  return y
                }
              })
              xyz = [...xyz, what]
              return el
            })
            states.data = states.data.filter((el) => {
              return data.list.every((el2) => {
                const eel = el2
                return eel.accountgs !== el.account
              })
            })
            states.arrayModal.datacount = xyz.join("\r\n")
            states.arrayModal.active = true
          }else{
            states.data = states.data.filter((el) => {
              return list.every((el2) => {
                return el2 !== el.account
              })
            })
            states.checkTemp = []
            states.loading = false
            ShowMessage(e)
          }
        }
      }else{
        states.data = states.data.filter((el) => {
          return list.every((el2) => {
            return el2 !== el.account
          })
        })
        states.checkTemp = []
        states.loading = false
        ShowMessage(e)
      }
    }

    const delAllSql = async() => {
      const list = states.checkTemp
      const e = {
          active: true,
          message: "删除失败",
          color: "is-danger",
          newtime: 0,
        }
      if (list.length > 0 ) {
        states.loading = true
        const token = localStorage.getItem("token")
        const params = {
          'list': states.checkTemp
        }
        const data = await Fetch(Config.Api.delallsql, params, "DELETE", token)
        if (data.status == -1) {
          states.data = []
          states.total = 0
          states.page = []
          localStorage.removeItem("token")
        }else if (data.status == 0) {
          states.data = states.data.filter((el) => {
            return list.every((el2) => {
              return el2 !== el.account
            })
          })
          states.checkTemp = []
          states.loading = false
          const e = {
            active: true,
            message: data.message,
            color: "is-success",
            newtime: 0,
          }
          ShowMessage(e)
        }else{
          states.loading = false
          ShowMessage(e)
        }
      }else{
        states.loading = false
        ShowMessage(e)
      }
    }

    const closeArray = () => {
      location.reload()
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
      copyIt,
      gotoAccount,
      checkBox,
      checkall,
      delit,
      getall,
      ShowMessage,
      closeArray,
      delToSql,
      delAllSql,
      sortExptime
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
