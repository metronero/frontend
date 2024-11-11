<script setup lang="js">
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const toast = useToast();
const { push } = useRouter();
const completeOn = ref(1);
const expireAfter = ref('1h');
const fiatCurrency = ref('EUR');

const completeOnOptions = [
    { name: '0-conf', code: 0 },
    { name: '1-conf', code: 1 },
    { name: '10-conf', code: 10 }
];

const fiatOptions = [
    { name: 'Euro', code: 'EUR' },
    { name: 'US Dollar', code: 'USD' }
];

const expireOptions = [
    { name: '5 minutes', code: '5m' },
    { name: '30 minutes', code: '30m' },
    { name: '1 hour', code: '1h' },
    { name: '6 hours', code: '6h' },
    { name: '12 hours', code: '12h' },
    { name: '1 day', code: '1d' }
];

function getSettings() {
    axios
        .get(import.meta.env.VITE_API_BASE + '/merchant', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            completeOn.value = response.data.complete_on;
            expireAfter.value = response.data.expire_after;
            fiatCurrency.value = response.data.fiat_currency;
        })
        .catch((error) => {
            console.log(error);
            if (error.response?.status === 401) {
                toast.add({ severity: 'error', summary: 'Not logged in', detail: error.code, life: 3000 });
                push({ path: '/auth/login' });
            } else {
                toast.add({ severity: 'error', summary: 'Failed to fetch merchant settings', detail: error.response.data.message, life: 3000 });
            }
        });
}

function saveSettings() {
    axios
        .post(import.meta.env.VITE_API_BASE + '/merchant', { complete_on: completeOn.value, expire_after: expireAfter.value, fiat_currency: fiatCurrency.value }, { withCredentials: true })
        .then(() => {
            toast.add({ severity: 'success', summary: 'Saved merchant settings', life: 3000 });
        })
        .catch((error) => {
            console.log(error);
            if (error.response?.status === 401) {
                toast.add({ severity: 'error', summary: 'Not logged in', detail: error.code, life: 3000 });
                push({ path: '/auth/login' });
            } else {
                toast.add({ severity: 'error', summary: 'Failed to save merchant settings', detail: error.response.data.message, life: 3000 });
            }
        });
}

onMounted(() => {
    getSettings();
});
</script>

<template>
    <div>
        <div class="card">
            <h1 class="font-semibold text-2xl mb-5">Account Settings</h1>
            <div class="mb-3">
                <p class="mb-2"><span class="font-bold">Complete on:</span> After how many confirmations this invoice should be marked as complete.</p>
                <Select v-model="completeOn" :options="completeOnOptions" optionLabel="name" optionValue="code" checkmark />
            </div>
            <div class="mb-3">
                <p class="mb-2"><span class="font-bold">Expire after:</span> If not paid, after how long should this invoice expire.</p>
                <Select v-model="expireAfter" :options="expireOptions" optionLabel="name" optionValue="code" checkmark />
            </div>
            <div class="mb-3">
                <p class="mb-2"><span class="font-bold">Default fiat currency:</span> Display values in this currency in PoS mode and payment pages.</p>
                <Select v-model="fiatCurrency" :options="fiatOptions" optionLabel="name" optionValue="code" checkmark />
            </div>
            <Button label="Save" class="mt-5" icon="pi pi-save" @click="saveSettings" />
        </div>
    </div>
</template>
