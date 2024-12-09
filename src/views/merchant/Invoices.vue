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
const invoiceDialog = ref(false);
const deleteInvoiceDialog = ref(false);
const invoice = ref({});
const selectedInvoices = ref([]);
const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});
const submitted = ref(false);
const apiBaseUrl = import.meta.env.VITE_API_BASE;

onMounted(() => {
    getInvoices();
});

function convertToPiconeros(xmrAmount) {
    return Math.round(xmrAmount * 1e12);
}

function piconerosToMonero(piconeros) {
    // Define the conversion factor: 1 Monero = 10^12 piconeros
    const PICONEROS_IN_MONERO = 1e12;

    // Convert piconeros to Monero (float)
    const monero = piconeros / PICONEROS_IN_MONERO;

    // Return the result as a fixed precision string, e.g., 5 decimal places
    return monero.toFixed(5); // Adjust precision as needed
}

function getInvoices() {
    axios
        .get(import.meta.env.VITE_API_BASE + '/merchant/invoice', { withCredentials: true })
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

function openNew() {
    invoice.value = { amount: null, order_id: '', accept_url: '', cancel_url: '', callback_url: '', extra_data: '' };
    submitted.value = false;
    invoiceDialog.value = true;
}

function hideDialog() {
    invoiceDialog.value = false;
    submitted.value = false;
}

function saveInvoice() {
    submitted.value = true;
    invoice.value.amount = Number(invoice.value.amount); // Convert amount to number

    if (invoice.value.amount && invoice.value.order_id) {
        axios
            .post(
                import.meta.env.VITE_API_BASE + '/merchant/invoice',
                {
                    amount: convertToPiconeros(invoice.value.amount),
                    order_id: invoice.value.order_id,
                    accept_url: invoice.value.accept_url,
                    cancel_url: invoice.value.cancel_url,
                    callback_url: invoice.value.callback_url,
                    extra_data: invoice.value.extra_data
                },
                { withCredentials: true }
            )
            .then((response) => {
                invoice.value.invoice_id = response.data.invoice_id;
                invoice.value.status = 'Pending';
                invoice.value.last_update = getCurrentDateAsISOString();
                invoices.value.unshift(invoice.value);
                toast.add({ severity: 'success', summary: 'Successful', detail: 'Invoice Created', life: 3000 });
                invoiceDialog.value = false;
                invoice.value = {};
            })
            .catch((error) => {
                toast.add({ severity: 'error', summary: 'Failed to create invoice', detail: error.message, life: 3000 });
            });
    }
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

function copyLinkToClipboard(invoiceId) {
    const link = `${apiBaseUrl}/p/page/${invoiceId}`;
    navigator.clipboard
        .writeText(link)
        .then(() => {
            toast.add({ severity: 'success', summary: 'Copied', detail: 'Invoice link copied to clipboard', life: 3000 });
        })
        .catch((error) => {
            toast.add({ severity: 'error', summary: 'Copy Failed', detail: error.message, life: 3000 });
        });
}

function getCurrentDateAsISOString() {
    const date = new Date();
    return date.toISOString();
}
function exportCSV() {
    dt.value.exportCSV();
}
</script>

<template>
    <div>
        <div class="card">
            <h1 class="font-semibold text-xl mb-5">Invoices</h1>

            <Toolbar class="mb-6">
                <template #start>
                    <Button label="New" icon="pi pi-plus" severity="secondary" class="mr-2" @click="openNew" />
                    <Button label="Delete" icon="pi pi-trash" severity="secondary" @click="confirmDeleteInvoice" :disabled="!selectedInvoices || !selectedInvoices.length" />
                </template>
                <template #end>
                    <Button label="Export" icon="pi pi-upload" severity="secondary" @click="exportCSV($event)" />
                </template>
            </Toolbar>

            <DataTable
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
                        <InputText v-model="filters['global'].value" placeholder="Search..." />
                    </div>
                </template>

                <Column selectionMode="multiple" style="width: 3rem" :exportable="false"></Column>
                <Column field="invoice_id" header="Invoice ID" sortable></Column>
                <Column field="amount" header="Amount" sortable>
                    <template #body="slotProps"> {{ piconerosToMonero(slotProps.data.amount) }} XMR</template></Column
                >
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
                    <template #body="slotProps"> <Button icon="pi pi-copy" label="Copy Link" @click="copyLinkToClipboard(slotProps.data.invoice_id)" /> </template>
                </Column>
            </DataTable>
        </div>

        <Dialog v-model:visible="invoiceDialog" :style="{ width: '450px' }" header="Invoice Details" :modal="true">
            <div class="flex flex-col gap-4">
                <div>
                    <label for="amount" class="block font-bold mb-2">Amount</label>
                    <InputText id="amount" v-model.trim="invoice.amount" type="number" required :invalid="submitted && !invoice.amount" />
                    <small v-if="submitted && !invoice.amount" class="text-red-500">Amount is required.</small>
                </div>
                <div>
                    <label for="order_id" class="block font-bold mb-2">Order ID</label>
                    <InputText id="order_id" v-model.trim="invoice.order_id" required :invalid="submitted && !invoice.order_id" />
                    <small v-if="submitted && !invoice.order_id" class="text-red-500">Order ID is required.</small>
                </div>
                <div>
                    <label for="accept_url" class="block font-bold mb-2">Accept URL</label>
                    <InputText id="accept_url" v-model.trim="invoice.accept_url" />
                </div>
                <div>
                    <label for="cancel_url" class="block font-bold mb-2">Cancel URL</label>
                    <InputText id="cancel_url" v-model.trim="invoice.cancel_url" />
                </div>
                <div>
                    <label for="callback_url" class="block font-bold mb-2">Callback URL</label>
                    <InputText id="callback_url" v-model.trim="invoice.callback_url" />
                </div>
                <div>
                    <label for="extra_data" class="block font-bold mb-2">Extra Data</label>
                    <InputText id="extra_data" v-model.trim="invoice.extra_data" />
                </div>
            </div>

            <template #footer>
                <Button label="Cancel" icon="pi pi-times" text @click="hideDialog" />
                <Button label="Save" icon="pi pi-check" @click="saveInvoice" />
            </template>
        </Dialog>

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
