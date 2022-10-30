<template>
  <nav class="navbar is-white">
    <div class="container">
      <div class="navbar-brand">
        <span class="navbar-item brand-text">
          <img :src="logo">Card Manage
        </span>
      </div>
      <div id="navMenu" class="navbar-menu">
        <div class="navbar-start">
          <router-link class="navbar-item" to="acclist">
            未提取帐号
          </router-link>
          <router-link class="navbar-item" to="acclistd">
            已提取帐号
          </router-link>
        </div>
      </div>
      <div class="navbar-end">
        <div class="navbar-item">
          <div class="field is-grouped">
            <div class="navbar-item has-dropdown is-hoverable">
              <span class="navbar-link">
                <span class="icon">
                  <i class="fa fa-user-circle-o"></i>
                </span>
                <span>
                用户中心
                </span>
              </span>
          
              <div class="navbar-dropdown">
                <div class="dropdown-item">
                  <p><span class="is-size-7">可以刷号! 努力賺錢.</span></p>
                </div>
                <hr class="dropdown-divider">
                <hr class="navbar-divider">
                <a class="navbar-item" @click="LogOut">
                  退出登陆
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>
<script>
import Config from '@/helper/config'
import setStorage from '@/helper/setStorage'
import { reactive, toRefs, defineComponent } from 'vue'
import { useRouter } from 'vue-router'
export default defineComponent({
  name: 'CardList',
  setup() {
    const router = useRouter()
    let states = reactive({
      logo: Config.images[2],
    })
    const LogOut =() =>{
      setStorage(false)
      router.push("/")
    }
    return {
      ...toRefs(states),
      LogOut
    }
  }
})
</script>

<style scoped>
nav.navbar {
  border-top: 4px solid #276cda;
}
</style>