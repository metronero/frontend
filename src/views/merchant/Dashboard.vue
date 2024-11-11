<script setup>
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useSessionStore } from '../../stores/session';

const sessionStore = useSessionStore();
const toast = useToast();
const { push } = useRouter();
const totalInvoices = ref(0);
const pendingInvoices = ref(0);
const sales = ref('0');
const isLoading = ref(false);
const isHealthy = ref(true);
const recentInvoices = ref();

function getDashboardInfo() {
    isLoading.value = true;
    axios
        .get(import.meta.env.VITE_API_BASE + '/health', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            isHealthy.value = response.data.healthy;
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
        .get(import.meta.env.VITE_API_BASE + '/merchant/invoice/recent', { withCredentials: true })
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
        .get(import.meta.env.VITE_API_BASE + '/merchant/stats', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            sales.value = piconerosToMonero(response.data.total_sales);
            pendingInvoices.value = response.data.pending;
            totalInvoices.value = response.data.total_invoices;
        })
        .catch((error) => {
            console.log(error);
            if (error.response.status == 401) {
                toast.add({ severity: 'error', summary: 'Not logged in', detail: error.code, life: 3000 });
                push({ path: '/auth/login' });
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
</script>

<template>
    <Toast />
    <div class="grid grid-cols-12 gap-8">
        <div class="col-span-12">
            <div class="card mb-0">
                <p class="font-semibold text-xl">Welcome {{ sessionStore.name }}!</p>
            </div>
        </div>
        <div class="col-span-12 lg:col-span-6 xl:col-span-3">
            <div class="card mb-0">
                <div class="flex justify-between mb-4">
                    <div>
                        <span class="block text-muted-color font-medium mb-4">Sales</span>
                        <div class="text-surface-900 dark:text-surface-0 font-medium text-xl">{{ sales }} XMR</div>
                    </div>
                    <div class="flex items-center justify-center bg-blue-100 dark:bg-blue-400/10 rounded-border" style="width: 2.5rem; height: 2.5rem">
                        <i class="pi pi-shopping-cart text-blue-500 !text-xl"></i>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-span-12 lg:col-span-6 xl:col-span-3">
            <div class="card mb-0">
                <div class="flex justify-between mb-4">
                    <div>
                        <span class="block text-muted-color font-medium mb-4">Invoices</span>
                        <div class="text-surface-900 dark:text-surface-0 font-medium text-xl">{{ totalInvoices }}</div>
                    </div>
                    <div class="flex items-center justify-center bg-orange-100 dark:bg-orange-400/10 rounded-border" style="width: 2.5rem; height: 2.5rem">
                        <i class="pi pi-dollar text-orange-500 !text-xl"></i>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-span-12 lg:col-span-6 xl:col-span-3">
            <div class="card mb-0">
                <div class="flex justify-between mb-4">
                    <div>
                        <span class="block text-muted-color font-medium mb-4">Pending</span>
                        <div class="text-surface-900 dark:text-surface-0 font-medium text-xl">{{ pendingInvoices }}</div>
                    </div>
                    <div class="flex items-center justify-center bg-cyan-100 dark:bg-cyan-400/10 rounded-border" style="width: 2.5rem; height: 2.5rem">
                        <i class="pi pi-users text-cyan-500 !text-xl"></i>
                    </div>
                </div>
            </div>
        </div>
        <div class="col-span-12 lg:col-span-6 xl:col-span-3">
            <div class="card mb-0">
                <div class="flex justify-between mb-4">
                    <div>
                        <span class="block text-muted-color font-medium mb-4">Instance</span>
                        <div class="text-surface-900 dark:text-surface-0 font-medium text-xl">{{ isHealthy ? 'Healthy' : 'Degraded' }}</div>
                    </div>
                    <div class="flex items-center justify-center bg-purple-100 dark:bg-purple-400/10 rounded-border" style="width: 2.5rem; height: 2.5rem">
                        <i class="pi pi-wave-pulse text-purple-500 !text-xl"></i>
                    </div>
                </div>
            </div>
        </div>

        <div class="col-span-12 xl:col-span-12">
            <div class="card">
                <div class="font-semibold text-xl mb-4">Recent Invoices</div>
                <DataTable v-if="recentInvoices" :value="recentInvoices" :rows="5" :paginator="true" responsiveLayout="scroll">
                    <Column field="invoice_id" header="Invoice ID" :sortable="true"></Column>
                    <Column header="Amount" :sortable="true">
                        <template #body="slotProps"> {{ piconerosToMonero(slotProps.data.amount) }} XMR </template>
                    </Column>
                    <Column field="order_id" header="Order ID" :sortable="true"></Column>
                    <Column field="status" header="Status" :sortable="true">
                        <template #body="slotProps">
                            <Tag v-if="slotProps.data.status === 'Completed'" value="Finished" severity="success" />
                            <Tag v-else-if="slotProps.data.status === 'Pending'" value="Pending" severity="warning" />
                            <Tag v-else-if="slotProps.data.status === 'Cancelled'" value="Cancelled" severity="danger" />
                            <Tag v-else value="Other" severity="info" />
                        </template>
                    </Column>
                    <Column field="last_update" header="Last Update" :sortable="true"></Column>
                </DataTable>
                <div class="flex items-center justify-center align-center" v-else>
                    <div class="text-center">
                        <p class="mb-4">Nothing yet. Create an invoice to get started.</p>
                        <Button as="router-link" to="/merchant/invoices">Go to Invoices</Button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
