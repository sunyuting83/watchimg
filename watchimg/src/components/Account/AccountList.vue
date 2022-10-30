<template>
  <div>
    <ManageHeader></ManageHeader>
    <div class="container">
      <nav class="columns">
        <div class="column">
          <div class="field is-grouped">
            <p class="control is-expanded">
              <input class="input is-small" type="text" v-model="SearchKey" placeholder="输入帐号-只能过滤当前页面数据" @input="Search">
            </p>
            <p class="control">
              <button class="button is-small is-info" :disabled="SearchKey.length > 0 ? false : true" @click="Clean">
                清空结果
              </button>
            </p>
          </div>
        </div>
      </nav>
      <div class="card events-card">
        <header class="card-header">
          <p class="card-header-title">
            {{title}}卡号列表
          </p>
        </header>
        <div class="card-content">
          <div class="content has-text-centered	min-heights" style="min-height: 11.3rem">
            <div class="com__box" v-if="loading" :style="loading? 'margin-top:5rem':''">
              <LoadIng></LoadIng>
            </div>
            <div v-else>
              <div v-if="data.length <= 0">
                <EmptyEd></EmptyEd>
              </div>
              <table class="table is-striped is-hoverable is-fullwidth is-narrow has-text-left" v-else>
                <thead>
                  <tr>
                    <td><label class="checkbox"><input type="checkbox" @click="checkall" />全选</label></td>
                    <td colspan="6"></td>
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
                    <td>炮台倍数</td>
                    <td>金币截图</td>
                    <td>日期</td>
                    <td>操作</td>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(item) in data" :key="item.id" @mouseover="setBackage(true,item.cover)" @mouseout="setBackage(false,'')">
                    <td><input type="checkbox" v-model="item.check" @click="(e)=>checkBox(e,item.account)"></td>
                    <td>{{item.account}}</td>
                    <td>{{makeNumber(item.today)}}</td>
                    <td>{{makeNumber(item.yesterday)}}</td>
                    <td>{{item.multiple}}</td>
                    <td><img class="thisimg" :src="rootUrl + item.cover" /></td>
                    <td><FormaTime :DateTime="item.datetime"></FormaTime></td>
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
        </div>
      </div>

      <PaginAtion v-if="data.length > 0" :total="total" :number="100" :GetData="GetData"></PaginAtion>
    </div>
    <NotIfication
      :showData="openerr">
    </NotIfication>
    <RenewalCard
      :showData="openModal"
      :ShowMessage="ShowMessage"></RenewalCard>
    <div class="is-img" :style="{top:hoverTop+'px'}" v-if="imghover"><img :src="rootUrl + currentImg" /></div>
  </div>
</template>
<script>
import { reactive, toRefs, defineComponent, onBeforeMount } from 'vue'
import { useRouter } from 'vue-router'
import ManageHeader from '@/components/Other/Header'
import LoadIng from '@/components/Other/Loading'
import EmptyEd from '@/components/Other/Empty'
import NotIfication from "@/components/Other/Notification"
import PaginAtion from '@/components/Other/PaginAtion'
import FormaTime from '@/components/Other/FormaTime'
import RenewalCard from '@/components/Other/Renewal'


import Fetch from '@/helper/fetch'
import Config from '@/helper/config'
export default defineComponent({
  name: 'AccList',
  components: { ManageHeader, LoadIng, EmptyEd, NotIfication, PaginAtion, FormaTime, RenewalCard },
  setup() {
    const router = useRouter()
    let states = reactive({
      loading: false,
      temp: [],
      checkTemp: [],
      data: [],
      total: 0,
      rootUrl: Config.RootU,
      path:router.currentRoute.value.path,
      page: [],
      title: "",
      openerr: {
        active: false,
        message: "",
        color: "",
        newtime: 0,
      },
      SearchKey: "",
      openModal:{
        active: false,
        account: {
          password: "",
          accountgs: "",
        },
        atest: "",
      },
      currentImg: '',
      hover: false,
      imghover: false,
      status: 0,
      ok: false,
      checkbos:[],
      hoverTop: 0,
      listcount: "",
      haslist: false,
    })
    onBeforeMount(async()=>{
      states.loading = true
      states.title = "未取"
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
    const GetData = async(page = 1) => {
      states.loading = true
      let url = Config.Api.getdata
      const token = localStorage.getItem("token")
      const d = await Fetch(url, {page:page}, 'GET', token)
      states.loading = true
      if (d.status == 1) {
        states.data = d.data.map((e)=>{
          e.check = false
          return e
        })
        states.temp = d.data
        states.total = d.total
        states.loading = false
      }else{
        states.data = []
        states.total = 0
        states.page = []
        states.loading = false
        localStorage.removeItem("token")
        states.hidden = false
      }
    }
    const Search = () => {
      const search = states.SearchKey.toUpperCase()
      if (search) {
        states.data = states.temp.filter((data) => {
          return Object.keys(data).some(() => {
            return String(data['account']).indexOf(search) !== -1
          })
        })
      }else {
        states.data = states.temp
      }
    }
    const Clean = () => {
      states.SearchKey = ""
      states.data = states.temp
    }
    const ShowMessage = (e) => {
      states.openerr = e
    }
    const showModel = (e) => {
      states.openModal.card = e
      states.openModal.active = true
    }

    const makeNumber = (n) =>{
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
    }

    const setBackage = (hover,img) => {
      states.currentImg = img
      states.imghover = hover
      let  scrollTop = document.documentElement.scrollTop || window.pageYOffset || document.body.scrollTop;
      if (scrollTop <= 61) scrollTop = 62
      states.hoverTop = scrollTop
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
        states.status = 0
        states.page = []
        localStorage.removeItem("token")
        states.hidden = false
      }else if (data.status == 1) {
        ShowMessage(e)
      }else {
        if (data.data !== "{}") {
          const x = JSON.parse(data.data)
          let d = ""
          states.data = states.data.filter((e) => {
            d = `${x.accountgs}----${x.password}----${makeNumber(e.today)}----${e.multiple}倍砲台`
            return e.account !== account
          })
          states.total = states.total - 1
          states.openModal.active = true
          states.openModal.account = x
          states.openModal.atest = d
          states.checkTemp = []
        }else{
          ShowMessage(e)
        }
      }
    }

    return {
      ...toRefs(states),
      GetData,
      Search,
      Clean,
      ShowMessage,
      showModel,
      makeNumber,
      setBackage,
      checkBox,
      checkall,
      delit
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
