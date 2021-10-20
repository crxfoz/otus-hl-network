import * as Vue from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createStore } from 'vuex'

import Login from './components/Login.vue'
import Registration from './components/Registration.vue'
import Profile from './components/Profile.vue'
import Users from './components/Users.vue'
import Friends from './components/Friends.vue'
import ProfileUpdate from './components/ProfileUpdate.vue'
import ProfileMy from './components/ProfileMy.vue'

import App from './App.vue'


const routes = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
    },
    {
        path: '/registration',
        name: 'Registration',
        component: Registration,
    },
    {
        path: '/me',
        name: 'ProfileMy',
        component: ProfileMy,
    },
    {
        path: '/profile/:id',
        name: 'Profile',
        component: Profile,
    },
    {
        path: '/update',
        name: 'ProfileUpdate',
        component: ProfileUpdate,
    },
    {
        path: '/users',
        name: 'Users',
        component: Users,
    },
    {
        path: '/friends',
        name: 'Friends',
        component: Friends,
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

const store = createStore({
    state: {
        authorized: false,
        user_id: 0,
        username: "",
        token: "",
    },
    mutations: {
        auth(state, payload) {
            state.authorized = payload.authorized;
            state.user_id = payload.user_id;
            state.username = payload.username;
            state.token = payload.token;
        }
    }
})

export default router

const app = Vue.createApp(App)
app.use(router)
app.use(store)
app.mount('#app')
