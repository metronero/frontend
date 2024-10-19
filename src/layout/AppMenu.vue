<script setup>
import axios from 'axios';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useSessionStore } from '../stores/session';
import AppMenuItem from './AppMenuItem.vue';

const { push } = useRouter();
const sessionStore = useSessionStore();

function logout() {
    axios
        .post(import.meta.env.VITE_API_BASE + '/logout', null, { withCredentials: true })
        .catch((error) => {
            console.log(error);
        })
        .finally(() => {
            localStorage.clear();
            document.cookie = 'MNSession=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
            push({ path: '/auth/login' });
        });
}

const adminMenu = ref([
    {
        label: 'Home',
        items: [
            { label: 'Dashboard', icon: 'pi pi-fw pi-home', to: '/' },
            { label: 'Invoices', icon: 'pi pi-fw pi-shopping-cart', to: '/admin/invoices' },
            { label: 'Withdrawals', icon: 'pi pi-fw pi-dollar', to: '/admin/withdrawals' },
            { label: 'Merchants', icon: 'pi pi-fw pi-users', to: '/admin/merchants' }
        ]
    },
    {
        label: 'Instance',
        items: [{ label: 'Info', icon: 'pi pi-fw pi-info', to: '/admin/instance' }]
    },
    {
        label: 'Account',
        items: [
            { label: 'Settings', icon: 'pi pi-fw pi-cog', to: '/admin/settings' },
            {
                label: 'Log out',
                icon: 'pi pi-fw pi-sign-out',
                command: () => {
                    logout();
                }
            }
        ]
    },
    {
        label: 'Support',
        items: [
            { label: 'Documentation', icon: 'pi pi-fw pi-book', url: 'https://metronero.moneropay.eu', target: '_blank' },
            { label: 'Source code', icon: 'pi pi-fw pi-github', url: 'https://github.com/metronero', target: '_blank' }
        ]
    }
]);

const merchantMenu = ref([
    {
        label: 'Home',
        items: [
            { label: 'Dashboard', icon: 'pi pi-fw pi-home', to: '/' },
            { label: 'Invoices', icon: 'pi pi-fw pi-shopping-cart', to: '/merchant/invoices' }
        ]
    },
    {
        label: 'Account',
        items: [
            { label: 'Settings', icon: 'pi pi-fw pi-cog', to: '/merchant/settings' },
            {
                label: 'Log out',
                icon: 'pi pi-fw pi-sign-out',
                command: () => {
                    logout();
                }
            }
        ]
    },
    {
        label: 'Support',
        items: [
            { label: 'Documentation', icon: 'pi pi-fw pi-book', url: 'https://metronero.moneropay.eu', target: '_blank' },
            { label: 'Source code', icon: 'pi pi-fw pi-github', url: 'https://github.com/metronero', target: '_blank' }
        ]
    }
]);
</script>

<template>
    <ul v-if="sessionStore.role == 'admin'" class="layout-menu">
        <template v-for="(item, i) in adminMenu" :key="item">
            <app-menu-item v-if="!item.separator" :item="item" :index="i"></app-menu-item>
            <li v-if="item.separator" class="menu-separator"></li>
        </template>
    </ul>
    <ul v-else class="layout-menu">
        <template v-for="(item, i) in merchantMenu" :key="item">
            <app-menu-item v-if="!item.separator" :item="item" :index="i"></app-menu-item>
            <li v-if="item.separator" class="menu-separator"></li>
        </template>
    </ul>
</template>

<style lang="scss" scoped></style>
