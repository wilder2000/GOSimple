import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import LoginView from '../views/LoginView.vue'
import MainLayout from '../views/MainLayout.vue'
import DashboardView from '../views/DashboardView.vue'
import UserList from '../views/UserList.vue'
import UserForm from '../views/UserForm.vue'
import RoleList from '../views/RoleList.vue'
import RoleForm from '../views/RoleForm.vue'
import GroupList from '../views/GroupList.vue'
import GroupForm from '../views/GroupForm.vue'
import DepartmentList from '../views/DepartmentList.vue'
import OperatorList from '../views/OperatorList.vue'
import UrlMappingList from '../views/UrlMappingList.vue'
import LogList from '../views/LogList.vue'

const routes = [
  { path: '/login', name: 'Login', component: LoginView },
  {
    path: '/',
    component: MainLayout,
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'Dashboard', component: DashboardView },
      { path: 'users', name: 'Users', component: UserList },
      { path: 'users/new', name: 'UserNew', component: UserForm },
      { path: 'users/:id/edit', name: 'UserEdit', component: UserForm, props: true },
      { path: 'roles', name: 'Roles', component: RoleList },
      { path: 'roles/new', name: 'RoleNew', component: RoleForm },
      { path: 'roles/:id/edit', name: 'RoleEdit', component: RoleForm, props: true },
      { path: 'groups', name: 'Groups', component: GroupList },
      { path: 'groups/new', name: 'GroupNew', component: GroupForm },
      { path: 'groups/:id/edit', name: 'GroupEdit', component: GroupForm, props: true },
      { path: 'departments', name: 'Departments', component: DepartmentList },
      { path: 'operators', name: 'Operators', component: OperatorList },
      { path: 'urlmappings', name: 'UrlMappings', component: UrlMappingList },
      { path: 'logs', name: 'Logs', component: LogList },
    ],
  },
]

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes,
})

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()
  if (to.name !== 'Login' && !auth.isLoggedIn()) {
    next({ name: 'Login' })
  } else {
    next()
  }
})

export default router
