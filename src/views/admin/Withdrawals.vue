<script setup>
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { onMounted, ref } from 'vue';

const withdrawalAmount = ref(0.0);
const withdrawalAddress = ref('');
const sweepAll = ref(false);
const toast = useToast();
const dt = ref();
const allWithdrawals = ref([]);

// Fetch withdrawals on mount
onMounted(() => {
    getWithdrawals();
});

function piconerosToMonero(piconeros) {
    // Define the conversion factor: 1 Monero = 10^12 piconeros
    const PICONEROS_IN_MONERO = 1e12;

    // Convert piconeros to Monero (float)
    const monero = piconeros / PICONEROS_IN_MONERO;

    // Return the result as a fixed precision string, e.g., 5 decimal places
    return monero.toFixed(5); // Adjust precision as needed
}

// Function to get all withdrawals and display in history
function getWithdrawals() {
    axios
        .get(import.meta.env.VITE_API_BASE + '/admin/withdrawal', { withCredentials: true })
        .then((response) => {
            allWithdrawals.value = response.data;
        })
        .catch((error) => {
            if (error.response) {
                toast.add({ severity: 'error', summary: 'Failed to get withdrawals', detail: error.response.data.message, life: 3000 });
            } else {
                toast.add({ severity: 'error', summary: error.message, detail: error.code, life: 3000 });
            }
        });
}

// Export CSV for withdrawals history
function exportCSV() {
    dt.value.exportCSV();
}

// Submit a withdrawal request
async function submitWithdrawal() {
    // Validate address is present
    if (!withdrawalAddress.value) {
        toast.add({ severity: 'warn', summary: 'Address Required', detail: 'Please enter a receiving address.', life: 3000 });
        return;
    }

    // Validate amount if sweepAll is false
    if (!sweepAll.value && !withdrawalAmount.value) {
        toast.add({ severity: 'warn', summary: 'Amount Required', detail: 'Please enter an amount or select Sweep All.', life: 3000 });
        return;
    }

    // Convert XMR amount to piconeros if sweepAll is false
    const amountInPiconeros = sweepAll.value ? undefined : Math.floor(withdrawalAmount.value * 1e12);

    // Construct request payload
    const requestPayload = {
        address: withdrawalAddress.value,
        amount: amountInPiconeros,
        sweep_all: sweepAll.value
    };

    try {
        const response = await axios.post(import.meta.env.VITE_API_BASE + '/admin/withdrawal', requestPayload, { withCredentials: true });

        // Show success toast with response data
        toast.add({
            severity: 'success',
            summary: 'Withdrawal Successful',
            detail: `Transaction ID: ${response.data.txid}, Amount: ${response.data.amount}`,
            life: 5000
        });

        // Insert the new withdrawal into the DataTable's data
        allWithdrawals.value.unshift({
            withdrawal_id: response.data.withdrawal_id, // Use the ID from the response as the unique identifier
            amount: response.data.amount, // Amount in piconeros
            address: withdrawalAddress.value,
            date: new Date().toISOString() // Current date or actual date if available in response
        });

        // Reset form fields
        withdrawalAddress.value = '';
        withdrawalAmount.value = 0;
        sweepAll.value = false;
    } catch (error) {
        // Handle errors and show relevant message
        const message = error.response ? error.response.data.message : error.message;
        toast.add({ severity: 'error', summary: 'Withdrawal Failed', detail: message, life: 3000 });
    }
}
</script>

<template>
    <div>
        <div class="card">
            <h1 class="font-semibold text-2xl mb-5">Withdrawals</h1>

            <h2 class="font-semibold text-xl mb-5">Initiate New</h2>

            <div class="mb-10">
                <div class="mb-5 flex flex-col">
                    <label for="address" class="font-bold block mb-2">Enter receiving address</label>
                    <InputText class="mb-5" required v-model="withdrawalAddress" inputId="address" />

                    <label for="withdrawalAmount" class="font-bold block mb-2">Enter amount</label>
                    <div class="flex gap-5">
                        <!-- Amount input disabled if sweepAll is checked -->
                        <InputNumber v-if="!sweepAll" v-model="withdrawalAmount" inputId="withdrawalAmount" suffix=" XMR" :minFractionDigits="0" :maxFractionDigits="10" />
                        <InputNumber v-else v-model="withdrawalAmount" inputId="withdrawalAmount" suffix=" XMR" disabled />

                        <!-- Sweep All checkbox -->
                        <Checkbox size="large" v-model="sweepAll" inputId="sweepAll" binary />
                        <label for="sweepAll">Sweep all</label>
                    </div>
                </div>

                <!-- Submit button calls submitWithdrawal -->
                <Button type="submit" label="Submit" @click="submitWithdrawal" />
            </div>

            <h2 class="font-semibold text-xl mb-5">History</h2>

            <Toolbar class="mb-6">
                <template #end>
                    <Button label="Export" icon="pi pi-upload" severity="secondary" @click="exportCSV($event)" />
                </template>
            </Toolbar>

            <DataTable v-if="allWithdrawals" :value="allWithdrawals" :rows="5" :paginator="true" responsiveLayout="scroll">
                <Column field="withdrawal_id" header="ID" :sortable="true"></Column>
                <Column field="amount" header="Amount" :sortable="true">
                    <template #body="slotProps"> {{ piconerosToMonero(slotProps.data.amount) }} XMR</template>
                </Column>
                <Column field="address" header="Address" :sortable="true"></Column>
                <Column field="date" header="Date" :sortable="true"></Column>
            </DataTable>

            <!-- Display message when there are no withdrawals -->
            <div class="flex items-center justify-center align-center" v-else>
                <div class="text-center">
                    <p>Nothing yet. Withdraw some funds to get started.</p>
                </div>
            </div>
        </div>
    </div>
</template>
