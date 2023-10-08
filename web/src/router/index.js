// Composables
import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    path: "/",
    component: () => import("@/layouts/default/Default.vue"),
    children: [
      {
        path: "/home",
        name: "Home",
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () =>
          import(/* webpackChunkName: "home" */ "@/views/Home.vue"),
      },
      {
        path: "/posts",
        name: "Posts",
        component: () => import("@/views/Posts.vue"),
      },
      {
        path: "/account",
        name: "Account",
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () =>
          import(/* webpackChunkName: "account" */ "@/views/Account.vue"),
      },
      //chats
      {
        path: "/chats",
        name: "Chats",
        component: () => import("@/views/Chats.vue"),
      },
      {
        path: "/bots",
        name: "Bots",
        component: () => import("@/views/Bots.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
