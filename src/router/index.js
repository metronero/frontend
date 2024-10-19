import AppLayout from '@/layout/AppLayout.vue';
import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            component: AppLayout,
            children: [
                {
                    path: '/',
                    redirect: '/auth/login'
                },
                {
                    path: '/admin/dashboard',
                    name: 'admin-dashboard',
                    component: () => import('@/views/admin/Dashboard.vue')
                },
                {
                    path: '/admin/invoices',
                    name: 'admin-invoices',
                    component: () => import('@/views/admin/Invoices.vue')
                },
                {
                    path: '/admin/withdrawals',
                    name: 'admin-withdrawals',
                    component: () => import('@/views/admin/Withdrawals.vue')
                },
                {
                    path: '/admin/merchants',
                    name: 'admin-merchants',
                    component: () => import('@/views/admin/Merchants.vue')
                },
                {
                    path: '/admin/instance',
                    name: 'admin-instance',
                    component: () => import('@/views/admin/Instance.vue')
                },
                {
                    path: '/admin/settings',
                    name: 'admin-settings',
                    component: () => import('@/views/admin/Settings.vue')
                },
                {
                    path: '/merchant/dashboard',
                    name: 'merchant-dashboard',
                    component: () => import('@/views/merchant/Dashboard.vue')
                },
                {
                    path: '/merchant/invoices',
                    name: 'merchant-invoices',
                    component: () => import('@/views/merchant/Invoices.vue')
                },
                {
                    path: '/merchant/settings',
                    name: 'merchant-settings',
                    component: () => import('@/views/merchant/Settings.vue')
                }
            ]
        },
        {
            path: '/auth/login',
            name: 'login',
            component: () => import('@/views/auth/Login.vue')
        }
    ]
});

export default router;
