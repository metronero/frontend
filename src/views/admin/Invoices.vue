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
                <Column field="invoice_id" header="Invoice ID" sortable style="min-width: 14rem"></Column>
                <Column field="amount" header="Amount" sortable style="min-width: 12rem"></Column>
                <Column field="order_id" header="Order ID" sortable style="min-width: 12rem"></Column>
                <Column field="status" header="Status" sortable style="min-width: 10rem">
                    <template #body="slotProps">
                        <Tag v-if="slotProps.data.status === 'finished'" value="Finished" severity="success" />
                        <Tag v-else-if="slotProps.data.status === 'Pending'" value="Pending" severity="warning" />
                        <Tag v-else-if="slotProps.data.status === 'Cancelled'" value="Cancelled" severity="danger" />
                        <Tag v-else value="Other" severity="info" />
                    </template>
                </Column>
                <Column field="last_update" header="Last Update" sortable style="min-width: 12rem"></Column>
                <Column header="Page Link" :exportable="false" style="min-width: 10rem">
                    <template #body="slotProps">
                        <Button icon="pi pi-link" outlined rounded severity="success" as="a" :href="apiBaseUrl + '/p/' + slotProps.data.invoice_id" />
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
