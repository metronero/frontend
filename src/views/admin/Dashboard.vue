<script setup>
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const toast = useToast();
const { push } = useRouter();
const totalInvoices = ref(0);
const pendingInvoices = ref(0);
const totalBalance = ref('0');
const unlockedBalance = ref('0');
const totalMerchants = ref(0);
const activeMerchants = ref(0);
const instanceHealth = ref(true);
const instanceVersion = ref('v0.0.0');
const isLoading = ref(true);

function getDashboardInfo() {
    isLoading.value = true;
    axios
        .get(import.meta.env.VITE_API_BASE + '/admin/invoice/count', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            totalInvoices.value = response.data.count;
            pendingInvoices.value = response.data.pending;
        })
        .catch((error) => {
            console.log(error);
            if (error.response.status == 401) {
                toast.add({ severity: 'error', summary: 'Not logged in', detail: error.code, life: 3000 });
                push({ path: '/auth/login' });
            } else {
                toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
            }
        });

    axios
        .get(import.meta.env.VITE_API_BASE + '/admin/balance', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            totalBalance.value = piconerosToMonero(response.data.total);
            unlockedBalance.value = piconerosToMonero(response.data.unlocked);
        })
        .catch((error) => {
            console.log(error);
            if (error.response) {
                toast.add({ severity: 'error', summary: 'Failed to fetch balance', detail: error.response.data.message, life: 3000 });
            } else {
                toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
            }
        });

    axios
        .get(import.meta.env.VITE_API_BASE + '/admin/merchant/count', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            totalMerchants.value = response.data.count;
            activeMerchants.value = response.data.active;
        })
        .catch((error) => {
            console.log(error);
            if (error.response) {
                toast.add({ severity: 'error', summary: 'Failed to fetch merchant count', detail: error.response.data.message, life: 3000 });
            } else {
                toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
            }
        });

    axios
        .get(import.meta.env.VITE_API_BASE + '/admin/instance', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            instanceHealth.value = response.data.healthy;
            instanceVersion.value = response.data.version;
        })
        .catch((error) => {
            console.log(error);
            if (error.response) {
                toast.add({ severity: 'error', summary: 'Failed to fetch instance info', detail: error.response.data.message, life: 3000 });
            } else {
                toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
            }
        })
        .finally(() => (isLoading.value = false));
}

function piconerosToMonero(piconeros) {
    // Define the conversion factor: 1 Monero = 10^12 piconeros
    const PICONEROS_IN_MONERO = 1e12;

    // Convert piconeros to Monero (float)
    const monero = piconeros / PICONEROS_IN_MONERO;

    // Return the result as a fixed precision string, e.g., 5 decimal places
    return monero.toFixed(5); // Adjust precision as needed
}

onMounted(() => {
    getDashboardInfo();
});

const formatCurrency = (value) => {
    return value.toLocaleString('en-US', { style: 'currency', currency: 'USD' });
};
</script>

<template>
    <Toast />
    <div class="grid grid-cols-12 gap-8">
        <div class="col-span-12">
            <div class="card mb-0">
                <p class="font-semibold text-xl">Welcome admin!</p>
            </div>
        </div>
        <div class="col-span-12 lg:col-span-6 xl:col-span-3">
            <div class="card mb-0">
                <div class="flex justify-between mb-4">
                    <div>
                        <span class="block text-muted-color font-medium mb-4">Invoices</span>
                        <div class="text-surface-900 dark:text-surface-0 font-medium text-xl">{{ totalInvoices }}</div>
                    </div>
                    <div class="flex items-center justify-center bg-blue-100 dark:bg-blue-400/10 rounded-border" style="width: 2.5rem; height: 2.5rem">
                        <i class="pi pi-shopping-cart text-blue-500 !text-xl"></i>
                    </div>
                </div>
                <span class="text-primary font-medium">{{ pendingInvoices }} </span>
                <span class="text-muted-color"> pending</span>
            </div>
        </div>
        <div class="col-span-12 lg:col-span-6 xl:col-span-3">
            <div class="card mb-0">
                <div class="flex justify-between mb-4">
                    <div>
                        <span class="block text-muted-color font-medium mb-4">Balance</span>
                        <div class="text-surface-900 dark:text-surface-0 font-medium text-xl">{{ totalBalance }} XMR</div>
                    </div>
                    <div class="flex items-center justify-center bg-orange-100 dark:bg-orange-400/10 rounded-border" style="width: 2.5rem; height: 2.5rem">
                        <i class="pi pi-dollar text-orange-500 !text-xl"></i>
                    </div>
                </div>
                <span class="text-muted-color">Unlocked: </span>
                <span class="text-primary font-medium">{{ unlockedBalance }} XMR</span>
            </div>
        </div>
        <div class="col-span-12 lg:col-span-6 xl:col-span-3">
            <div class="card mb-0">
                <div class="flex justify-between mb-4">
                    <div>
                        <span class="block text-muted-color font-medium mb-4">Merchants</span>
                        <div class="text-surface-900 dark:text-surface-0 font-medium text-xl">{{ totalMerchants }}</div>
                    </div>
                    <div class="flex items-center justify-center bg-cyan-100 dark:bg-cyan-400/10 rounded-border" style="width: 2.5rem; height: 2.5rem">
                        <i class="pi pi-users text-cyan-500 !text-xl"></i>
                    </div>
                </div>
                <span class="text-primary font-medium">{{ activeMerchants }}</span>
                <span class="text-muted-color"> active</span>
            </div>
        </div>
        <div class="col-span-12 lg:col-span-6 xl:col-span-3">
            <div class="card mb-0">
                <div class="flex justify-between mb-4">
                    <div>
                        <span class="block text-muted-color font-medium mb-4">Instance</span>
                        <div class="text-surface-900 dark:text-surface-0 font-medium text-xl">{{ instanceHealth ? 'Healthy' : 'Degraded' }}</div>
                    </div>
                    <div class="flex items-center justify-center bg-purple-100 dark:bg-purple-400/10 rounded-border" style="width: 2.5rem; height: 2.5rem">
                        <i class="pi pi-wave-pulse text-purple-500 !text-xl"></i>
                    </div>
                </div>
                <span class="text-muted-color">Version: </span>
                <span class="text-primary font-medium">{{ instanceVersion }}</span>
            </div>
        </div>

        <div class="col-span-12 xl:col-span-6">
            <div class="card">
                <div class="font-semibold text-xl mb-4">Recent Invoices</div>
                <DataTable :value="products" :rows="5" :paginator="true" responsiveLayout="scroll">
                    <Column field="name" header="Name" :sortable="true" style="width: 35%"></Column>
                    <Column field="price" header="Price" :sortable="true" style="width: 35%">
                        <template #body="slotProps">
                            {{ formatCurrency(slotProps.data.price) }}
                        </template>
                    </Column>
                </DataTable>
            </div>
        </div>
        <div class="col-span-12 xl:col-span-6">
            <div class="card">
                <div class="font-semibold text-xl mb-4">Recent Withdrawals</div>
                <DataTable :value="products" :rows="5" :paginator="true" responsiveLayout="scroll">
                    <Column field="name" header="Name" :sortable="true" style="width: 35%"></Column>
                    <Column field="price" header="Price" :sortable="true" style="width: 35%">
                        <template #body="slotProps">
                            {{ formatCurrency(slotProps.data.price) }}
                        </template>
                    </Column>
                </DataTable>
            </div>
        </div>
    </div>
</template>
