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
const recentInvoices = ref();
const recentWithdrawals = ref();

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
        .get(import.meta.env.VITE_API_BASE + '/admin/invoice/recent', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            recentInvoices.value = response.data;
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
        .get(import.meta.env.VITE_API_BASE + '/admin/withdrawal/recent', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            recentWithdrawals.value = response.data;
        })
        .catch((error) => {
            console.log(error);
            if (error.response) {
                toast.add({ severity: 'error', summary: 'Failed to fetch recent withdrawals', detail: error.response.data.message, life: 3000 });
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

function copyAddress(address) {
    navigator.clipboard
        .writeText(address)
        .then(() => {
            // Show toast notification
            toast.add({ severity: 'success', summary: 'Address Copied', detail: 'Address copied to clipboard', life: 2000 });
        })
        .catch((err) => {
            // Handle any errors with copying
            console.error('Failed to copy: ', err);
            toast.add({ severity: 'error', summary: 'Copy Failed', detail: 'Could not copy address', life: 2000 });
        });
}

function formatDate(dateString) {
    // Create a new Date object from the ISO string
    const date = new Date(dateString);

    // Extract components (year, month, day, hours, minutes)
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0'); // Month is zero-indexed
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');

    // Return the formatted date as "YYYY-MM-DD HH:MM"
    return `${year}-${month}-${day} ${hours}:${minutes}`;
}

onMounted(() => {
    getDashboardInfo();
});
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
                <DataTable v-if="recentInvoices" :value="recentInvoices" :rows="5" :paginator="true" responsiveLayout="scroll">
                    <Column field="merchant_name" header="Merchant" sortable></Column>
                    <Column header="Amount" :sortable="true">
                        <template #body="slotProps"> {{ piconerosToMonero(slotProps.data.amount) }} XMR</template>
                    </Column>
                    <Column field="order_id" header="Order ID" :sortable="true"></Column>
                    <Column field="status" header="Status" sortable>
                        <template #body="slotProps">
                            <Tag v-if="slotProps.data.status === 'Completed'" value="Finished" severity="success" />
                            <Tag v-else-if="slotProps.data.status === 'Pending'" value="Pending" severity="info" />
                            <Tag v-else-if="slotProps.data.status === 'Partial'" value="Partial" severity="warning" />
                            <Tag v-else-if="slotProps.data.status === 'Cancelled'" value="Cancelled" severity="danger" />
                            <Tag v-else-if="slotProps.data.status === 'Expired'" value="Expired" severity="danger" />
                            <Tag v-else value="Other" severity="info" />
                        </template>
                    </Column>
                    <Column field="last_update" header="Last Update" sortable>
                        <template #body="slotProps">
                            {{ formatDate(slotProps.data.last_update) }}
                        </template>
                    </Column>
                </DataTable>
                <div class="flex items-center justify-center align-center" v-else>
                    <div class="text-center">
                        <p class="mb-4">Nothing yet. Get a merchant to create an invoice.</p>
                        <Button as="router-link" to="/admin/merchants">Go to Merchants</Button>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-span-12 xl:col-span-6">
            <div class="card">
                <div class="font-semibold text-xl mb-4">Recent Withdrawals</div>
                <DataTable v-if="recentWithdrawals" :value="recentWithdrawals" :rows="5" :paginator="true" responsiveLayout="scroll">
                    <Column field="amount" header="Amount" :sortable="true">
                        <template #body="slotProps"> {{ piconerosToMonero(slotProps.data.amount) }} XMR </template>
                    </Column>
                    <Column field="address" header="Address" :sortable="true">
                        <template #body="slotProps"> {{ slotProps.data.address.slice(0, 10) }}... <Button icon="pi pi-copy" text @click="copyAddress(slotProps.data.address)" /> </template>
                    </Column>
                    <Column field="date" header="Date" sortable>
                        <template #body="slotProps">
                            {{ formatDate(slotProps.data.date) }}
                        </template>
                    </Column>
                </DataTable>
                <div class="flex items-center justify-center align-center" v-else>
                    <div class="text-center">
                        <p class="mb-4">Nothing yet. Withdraw some funds to get started.</p>
                        <Button as="router-link" to="/admin/withdrawals">Go to Withdrawals</Button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
