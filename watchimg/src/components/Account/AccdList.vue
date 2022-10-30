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
        <div class="card events-card mt-5"  v-else  v-for="(item) in data" :key="item.date">
          <header class="card-header">
            <p class="card-header-title">
              {{title}}卡号列表-日期：{{item.date}}
            </p>
          </header>
          <div class="card-content">
            <div class="content has-text-centered	min-heights" style="min-height: 11.3rem">
              <table class="table is-striped is-hoverable is-fullwidth is-narrow has-text-left">
                <thead>
                  <tr>
                    <td>帐号</td>
                    <td>今日金币</td>
                    <td>炮台倍数</td>
                    <td>金币截图</td>
                    <td>提号人</td>
                    <td>日期</td>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(it) in item.data" :key="it.id" @mouseover="setBackage(true,it.cover)" @mouseout="setBackage(false,'')">
                    <td>{{it.account}}</td>
                    <td>{{makeNumber(it.today)}}</td>
                    <td>{{it.multiple}}</td>
                    <td><img class="thisimg" :src="rootUrl + it.cover" /></td>
                    <td>{{it.username}}</td>
                    <td><FormaTime :DateTime="it.updatetime"></FormaTime></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <PaginAtion v-if="data.length > 0" :total="total" :number="pageNumber" :GetData="makePageData"></PaginAtion>
    </div>
    <div class="is-img" :style="{top:hoverTop+'px'}" v-if="imghover"><img :src="rootUrl + currentImg" /></div>
  </div>
</template>
<script>
import { reactive, toRefs, defineComponent, onBeforeMount } from 'vue'
import { useRouter } from 'vue-router'

import ManageHeader from '@/components/Other/Header'
import LoadIng from '@/components/Other/Loading'
import EmptyEd from '@/components/Other/Empty'
import PaginAtion from '@/components/Other/PaginAtion'
import FormaTime from '@/components/Other/FormaTime'


import Fetch from '@/helper/fetch'
import Config from '@/helper/config'
export default defineComponent({
  name: 'AccdList',
  components: { ManageHeader, LoadIng, EmptyEd, PaginAtion, FormaTime },
  setup() {
    const router = useRouter()
    let states = reactive({
      loading: false,
      temp: [],
      data: [],
      SearchTemp:[],
      total: 0,
      rootUrl: Config.RootU,
      page: [],
      title: "",
      SearchKey: "",
      currentImg: '',
      imghover: false,
      hoverTop: 0,
      pageNumber: 3,
      currentPage: 1,
      SearchDatekey: ""
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
      let url = Config.Api.getodata
      const token = localStorage.getItem("token")
      const d = await Fetch(url, {}, 'GET', token)
      states.loading = true
      if (d.status == 1) {
        const NewData = makData(d.data)
        states.SearchTemp = d.data
        states.temp = NewData
        states.total = NewData.length
        makePageData(1)
      }else{
        states.data = []
        states.total = 0
        states.page = []
        states.loading = false
      }
    }
    const Search = () => {
      const search = states.SearchKey.toUpperCase()
      if (search) {
        const searchData = states.SearchTemp.filter((data) => {
          return Object.keys(data).some(() => {
            return String(data['account']).toUpperCase().indexOf(search) !== -1
          })
        })
        states.data = makData(searchData)
      }else {
        makePageData(states.currentPage)
      }
    }
    const Clean = () => {
      states.SearchKey = ""
      makePageData(states.currentPage)
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

    const setBackage = (hover,img) => {
      states.currentImg = img
      states.imghover = hover
      let  scrollTop = document.documentElement.scrollTop || window.pageYOffset || document.body.scrollTop;
      states.hoverTop = scrollTop + 60
    }
    const makePageData = (page) => {
      states.loading = true
      const offSet = page * states.pageNumber
      const limit = (page - 1) * states.pageNumber
      const d = states.temp
      const x = d.slice(limit, offSet)
      states.data = x
      states.loading = false
      states.currentPage = page
    }
    const SearchDate = () => {
      const search = states.SearchDatekey
      if (search) {
        states.data = states.temp.filter((data) => {
          return Object.keys(data).some(() => {
            return String(data['date']).indexOf(search) !== -1
          })
        })
      }else {
        makePageData(states.currentPage)
      }
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
