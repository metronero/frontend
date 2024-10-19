<script setup>
import FloatingConfigurator from '@/components/FloatingConfigurator.vue';
import axios from 'axios';
import ProgressSpinner from 'primevue/progressspinner';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useSessionStore } from '../../stores/session';

const sessionStore = useSessionStore();
const toast = useToast();
const { push } = useRouter();
const isLoading = ref(false);
const username = ref('');
const password = ref('');

function login() {
    isLoading.value = true;
    axios
        .post(import.meta.env.VITE_API_BASE + '/login', { username: username.value, password: password.value }, { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            localStorage.setItem('account_role', response.data.role);
            localStorage.setItem('account_name', username.value);
            sessionStore.setName(username.value);
            sessionStore.setRole(response.data.role);
            if (response.data.role == 'admin') {
                push({ path: '/admin/dashboard' });
            } else {
                push({ path: '/merchant/dashboard' });
            }
        })
        .catch((error) => {
            console.log(error);
            if (error.response) {
                toast.add({ severity: 'error', summary: 'Login Failed', detail: error.response.data.message, life: 3000 });
            } else {
                toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
            }
        })
        .finally(() => (isLoading.value = false));
}
</script>

<template>
    <Toast />
    <FloatingConfigurator />
    <div class="bg-surface-50 dark:bg-surface-950 flex items-center justify-center min-h-screen min-w-[100vw] overflow-hidden">
        <div class="flex flex-col items-center justify-center">
            <div style="border-radius: 56px; padding: 0.3rem; background: linear-gradient(180deg, var(--primary-color) 10%, rgba(33, 150, 243, 0) 30%)">
                <div class="w-full bg-surface-0 dark:bg-surface-900 py-20 px-8 sm:px-20" style="border-radius: 53px">
                    <div v-if="isLoading" class="flex flex-col items-center justify-center">
                        <ProgressSpinner />
                    </div>
                    <div v-else class="text-center mb-8">
                        <img class="mb-8 w-24 shrink-0 mx-auto" src="/metronero-logo.png" />
                        <div class="text-surface-900 dark:text-surface-0 text-3xl font-medium mb-4">Welcome to Metronero!</div>
                        <span class="text-muted-color font-medium">Sign in to continue</span>
                    </div>

                    <div>
                        <label for="username1" class="block text-surface-900 dark:text-surface-0 text-xl font-medium mb-2">Username</label>
                        <InputText id="username1" type="text" placeholder="Username" class="w-full md:w-[30rem] mb-8" v-model="username" />

                        <label for="password1" class="block text-surface-900 dark:text-surface-0 font-medium text-xl mb-2">Password</label>
                        <Password id="password1" v-model="password" placeholder="Password" :toggleMask="true" class="mb-4" fluid :feedback="false"></Password>

                        <Button label="Sign In" class="mt-8 w-full" @click="login()"></Button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.pi-eye {
    transform: scale(1.6);
    margin-right: 1rem;
}

.pi-eye-slash {
    transform: scale(1.6);
    margin-right: 1rem;
}
</style>
