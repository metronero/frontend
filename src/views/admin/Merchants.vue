<script setup lang="js">
import { FilterMatchMode } from '@primevue/core/api';
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

onMounted(() => {
    getMerchants();
});

const { push } = useRouter();
const toast = useToast();
const dt = ref();
const merchants = ref([]);
const merchantDialog = ref(false);
const deleteMerchantDialog = ref(false);
const deleteMerchantsDialog = ref(false);
const merchant = ref({});
const selectedMerchants = ref();
const filters = ref({
    global: { value: null, matchMode: FilterMatchMode.CONTAINS }
});
const submitted = ref(false);
const isEditMode = computed(() => !!merchant.value.id);

function getMerchants() {
    axios
        .get(import.meta.env.VITE_API_BASE + '/admin/merchant', { withCredentials: true })
        .then((response) => {
            console.log(response.data);
            merchants.value = response.data;
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
}

function openNew() {
    merchant.value = {};
    submitted.value = false;
    merchantDialog.value = true;
}

function hideDialog() {
    merchantDialog.value = false;
    submitted.value = false;
}

function saveMerchant() {
    submitted.value = true;

    // Check required fields based on whether it's an edit or create action
    if (merchant.value.name?.trim() && (isEditMode.value || merchant.value.password)) {
        if (isEditMode.value) {
            axios
                .post(import.meta.env.VITE_API_BASE + '/admin/merchant/' + merchant.value.id, { disabled: merchant.value.disabled }, { withCredentials: true })
                .then(() => {
                    merchants.value[findIndexById(merchant.value.id)] = merchant.value;
                    toast.add({ severity: 'success', summary: 'Successful', detail: 'Merchant Updated', life: 3000 });
                })
                .catch((error) => {
                    if (error.response) {
                        toast.add({ severity: 'error', summary: 'Failed to update merchant', detail: error.response.data.message, life: 3000 });
                    } else {
                        toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
                    }
                });
        } else {
            axios
                .post(import.meta.env.VITE_API_BASE + '/admin/register', { username: merchant.value.name, password: merchant.value.password, role: 'merchant' }, { withCredentials: true })
                .then((response) => {
                    merchant.value.id = response.data.id;
                    merchant.value.password = '';
                    merchants.value.push(merchant.value);
                    toast.add({ severity: 'success', summary: 'Successful', detail: 'Merchant Created', life: 3000 });
                    merchantDialog.value = false;
                    merchant.value = {};
                })
                .catch((error) => {
                    if (error.response) {
                        toast.add({ severity: 'error', summary: 'Failed to fetch merchants', detail: error.response.data.message, life: 3000 });
                    } else {
                        toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
                    }
                });
        }
        hideDialog();
    }
}

function editMerchant(m) {
    merchant.value = { ...m };
    merchantDialog.value = true;
}

function deleteMerchant() {
    merchants.value = merchants.value.filter((val) => val.id !== merchant.value.id);
    deleteMerchantDialog.value = false;
    merchant.value = {};
    toast.add({ severity: 'success', summary: 'Successful', detail: 'Merchant Deleted', life: 3000 });
}

function findIndexById(id) {
    let index = -1;
    for (let i = 0; i < merchants.value.length; i++) {
        if (merchants.value[i].id === id) {
            index = i;
            break;
        }
    }
    return index;
}

function exportCSV() {
    dt.value.exportCSV();
}

function confirmDeleteSelected() {
    deleteMerchantsDialog.value = true;
}

function deleteSelectedMerchants() {
    merchants.value = merchants.value.filter((val) => !selectedMerchants.value.includes(val));
    deleteMerchantsDialog.value = false;
    selectedMerchants.value = null;
    toast.add({ severity: 'success', summary: 'Successful', detail: 'Merchants Deleted', life: 3000 });
}
</script>

<template>
    <div>
        <div class="card">
            <Toolbar class="mb-6">
                <template #start>
                    <Button label="New" icon="pi pi-plus" severity="secondary" class="mr-2" @click="openNew" />
                    <Button label="Delete" icon="pi pi-trash" severity="secondary" @click="confirmDeleteSelected" :disabled="!selectedMerchants || !selectedMerchants.length" />
                </template>

                <template #end>
                    <Button label="Export" icon="pi pi-upload" severity="secondary" @click="exportCSV($event)" />
                </template>
            </Toolbar>

            <DataTable
                ref="dt"
                v-model:selection="selectedMerchants"
                :value="merchants"
                dataKey="id"
                :paginator="true"
                :rows="10"
                :filters="filters"
                paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
                :rowsPerPageOptions="[5, 10, 25]"
                currentPageReportTemplate="Showing {first} to {last} of {totalRecords} merchants"
            >
                <template #header>
                    <div class="flex flex-wrap gap-2 items-center justify-between">
                        <h4 class="m-0">Manage Merchants</h4>
                        <IconField>
                            <InputIcon>
                                <i class="pi pi-search" />
                            </InputIcon>
                            <InputText v-model="filters['global'].value" placeholder="Search..." />
                        </IconField>
                    </div>
                </template>

                <Column selectionMode="multiple" style="width: 3rem" :exportable="false"></Column>
                <Column field="id" header="ID" sortable style="min-width: 12rem"></Column>
                <Column field="name" header="Name" sortable style="min-width: 16rem"></Column>
                <Column field="disabled" header="Status" sortable style="min-width: 12rem">
                    <template #body="slotProps">
                        <Tag v-if="slotProps.data.disabled" value="Disabled" severity="error" />
                        <Tag v-else value="Active" severity="success" />
                    </template>
                </Column>
                <Column :exportable="false" style="min-width: 12rem">
                    <template #body="slotProps">
                        <Button icon="pi pi-pencil" class="mr-2" @click="editMerchant(slotProps.data)" label="Edit" />
                    </template>
                </Column>
            </DataTable>
        </div>

        <Dialog v-model:visible="merchantDialog" :style="{ width: '450px' }" header="Merchant Details" :modal="true">
            <div class="flex flex-col gap-6">
                <div v-if="!isEditMode">
                    <label for="name" class="block font-bold mb-3">Name</label>
                    <InputText id="name" v-model.trim="merchant.name" required="true" autofocus :invalid="submitted && !merchant.name" fluid />
                    <small v-if="submitted && !merchant.name" class="text-red-500">Name is required.</small>
                </div>
                <div v-if="!isEditMode">
                    <label for="password" class="block font-bold mb-3">Password</label>
                    <Password id="password" v-model.trim="merchant.password" autofocus :invalid="submitted && !merchant.password" fluid />
                    <small v-if="submitted && !merchant.password" class="text-red-500">Password is required.</small>
                </div>
                <div v-if="isEditMode">
                    <Checkbox v-model="merchant.disabled" binary inputId="disabled" class="mr-2" />
                    <label for="disabled">Disabled</label>
                </div>
            </div>

            <template #footer>
                <Button label="Cancel" icon="pi pi-times" text @click="hideDialog" />
                <Button label="Save" icon="pi pi-check" @click="saveMerchant" />
            </template>
        </Dialog>

        <Dialog v-model:visible="deleteMerchantDialog" :style="{ width: '450px' }" header="Confirm" :modal="true">
            <div class="flex items-center gap-4">
                <i class="pi pi-exclamation-triangle !text-3xl" />
                <span v-if="merchant"
                    >Are you sure you want to delete <b>{{ merchant.name }}</b
                    >?</span
                >
            </div>
            <template #footer>
                <Button label="No" icon="pi pi-times" text @click="deleteMerchantDialog = false" />
                <Button label="Yes" icon="pi pi-check" @click="deleteMerchant" />
            </template>
        </Dialog>

        <Dialog v-model:visible="deleteMerchantsDialog" :style="{ width: '450px' }" header="Confirm" :modal="true">
            <div class="flex items-center gap-4">
                <i class="pi pi-exclamation-triangle !text-3xl" />
                <span v-if="merchant">Are you sure you want to delete the selected merchants?</span>
            </div>
            <template #footer>
                <Button label="No" icon="pi pi-times" text @click="deleteMerchantsDialog = false" />
                <Button label="Yes" icon="pi pi-check" text @click="deleteSelectedMerchants" />
            </template>
        </Dialog>
    </div>
</template>
