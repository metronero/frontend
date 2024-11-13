<script setup lang="js">
import { FilterMatchMode } from '@primevue/core/api';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

const { push } = useRouter();
const toast = useToast();
const dt = ref();
const invoices = ref([]);
const deleteInvoiceDialog = ref(false);
const invoice = ref({});
const selectedInvoices = ref([]);
const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});
const apiBaseUrl = import.meta.env.VITE_API_BASE;

onMounted(() => {
    getInvoices();
});

function getInvoices() {
    axios
        .get(import.meta.env.VITE_API_BASE + '/admin/invoice', { withCredentials: true })
        .then((response) => {
            invoices.value = response.data;
        })
        .catch((error) => {
            if (error.response?.status === 401) {
                toast.add({ severity: 'error', summary: 'Not logged in', detail: error.code, life: 3000 });
                push({ path: '/auth/login' });
            } else {
                toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
            }
        });
}

function confirmDeleteInvoice(inv) {
    invoice.value = inv;
    deleteInvoiceDialog.value = true;
}

function deleteInvoice() {
    invoices.value = invoices.value.filter((val) => val.invoice_id !== invoice.value.invoice_id);
    deleteInvoiceDialog.value = false;
    toast.add({ severity: 'success', summary: 'Successful', detail: 'Invoice Deleted', life: 3000 });
}

function exportCSV() {
    dt.value.exportCSV();
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

function piconerosToMonero(piconeros) {
    // Define the conversion factor: 1 Monero = 10^12 piconeros
    const PICONEROS_IN_MONERO = 1e12;

    // Convert piconeros to Monero (float)
    const monero = piconeros / PICONEROS_IN_MONERO;

    // Return the result as a fixed precision string, e.g., 5 decimal places
    return monero.toFixed(5); // Adjust precision as needed
}
</script>

<template>
    <div>
        <div class="card">
            <Toolbar class="mb-6">
                <template #start>
                    <Button label="Delete" icon="pi pi-trash" severity="secondary" @click="confirmDeleteInvoice" :disabled="!selectedInvoices || !selectedInvoices.length" />
                </template>
                <template #end>
                    <Button label="Export" icon="pi pi-upload" severity="secondary" @click="exportCSV($event)" />
                </template>
            </Toolbar>

            <DataTable
                v-if="invoices"
                ref="dt"
                v-model:selection="selectedInvoices"
                :value="invoices"
                dataKey="invoice_id"
                :paginator="true"
                :rows="10"
                :filters="filters"
                paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
                :rowsPerPageOptions="[5, 10, 25]"
                currentPageReportTemplate="Showing {first} to {last} of {totalRecords} invoices"
            >
                <template #header>
                    <div class="flex justify-between">
                        <h4 class="m-0">Manage Invoices</h4>
                        <InputText v-model="filters['global'].value" placeholder="Search..." />
                    </div>
                </template>

                <Column selectionMode="multiple" style="width: 3rem" :exportable="false"></Column>
                <Column field="invoice_id" header="Invoice ID" sortable></Column>
                <Column field="merchant_name" header="Merchant" sortable></Column>
                <Column header="Amount" :sortable="true">
                    <template #body="slotProps"> {{ piconerosToMonero(slotProps.data.amount) }} XMR</template>
                </Column>
                <Column field="order_id" header="Order ID" sortable></Column>
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
                <Column header="Page Link" :exportable="false">
                    <template #body="slotProps">
                        <Button icon="pi pi-link" outlined rounded severity="success" as="a" :href="apiBaseUrl + '/p/page/' + slotProps.data.invoice_id" />
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

        <Dialog v-model:visible="deleteInvoiceDialog" :style="{ width: '450px' }" header="Confirm" :modal="true">
            <div class="flex items-center gap-4">
                <i class="pi pi-exclamation-triangle !text-3xl" />
                <span
                    >Are you sure you want to delete <b>{{ invoice.invoice_id }}</b
                    >?</span
                >
            </div>
            <template #footer>
                <Button label="No" icon="pi pi-times" text @click="deleteInvoiceDialog = false" />
                <Button label="Yes" icon="pi pi-check" @click="deleteInvoice" />
            </template>
        </Dialog>
    </div>
</template>
