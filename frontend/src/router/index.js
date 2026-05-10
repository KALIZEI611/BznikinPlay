import { createRouter, createWebHistory } from "vue-router";
import HomePage from "../views/HomePage.vue";
import MyRentals from "../views/MyRentals.vue";
import CheckoutPage from "../views/CheckoutPage.vue";
import ProfilePage from "../views/ProfilePage.vue";

const routes = [
  { path: "/", component: HomePage },
  { path: "/rentals", component: MyRentals, meta: { requiresAuth: true } },
  {
    path: "/checkout/:id?",
    component: CheckoutPage,
    name: "checkout",
    meta: { requiresAuth: true },
  },
  {
    path: "/profile",
    component: ProfilePage,
    name: "profile",
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem("auth_token");
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);

  if (requiresAuth && !token) {
    next("/");
  } else {
    next();
  }
});

export default router;
