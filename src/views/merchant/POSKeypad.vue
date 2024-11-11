<script setup lang="js">
import axios from 'axios';
import { useToast } from 'primevue/usetoast';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

// Define currencies
const selectedCurrency = ref('XMR');
const currencies = ref(['XMR', 'EUR', 'USD']);
const exchangeRate = ref(1); // Exchange rate for fiat to XMR conversion
const amount = ref(0); // Amount in the selected currency
const amountString = ref(''); // Track the amount as a string
const toast = useToast();
const { push } = useRouter();

// Invoice data
const invoice = ref({
    amount: 0,
    order_id: '',
    accept_url: window.location.href,
    cancel_url: window.location.href
});

// Convert amount in XMR to piconeros (1 XMR = 1e12 piconeros)
function convertToPiconeros(xmrAmount) {
    return Math.round(xmrAmount * 1e12);
}

function fetchExchangeRate() {
    if (selectedCurrency.value === 'XMR') {
        exchangeRate.value = 1; // XMR to XMR rate is 1
    } else {
        // Use a synchronous approach with a promise
        return axios
            .get(`${import.meta.env.VITE_API_BASE}/merchant/fiatrate/${selectedCurrency.value}`, { withCredentials: true })
            .then((response) => {
                exchangeRate.value = response.data.price;
            })
            .catch((error) => {
                if (error.response?.status === 401) {
                    toast.add({ severity: 'error', summary: 'Not logged in', detail: error.code, life: 3000 });
                    push({ path: '/auth/login' });
                } else {
                    console.error(`Failed to fetch exchange rate for ${selectedCurrency.value}:`, error);
                    toast.add({ severity: 'error', summary: 'Exchange Rate Error', detail: 'Could not fetch exchange rate.', life: 3000 });
                }
            });
    }
}

// Generate a random order_id for the invoice
function generateOrderId() {
    return 'POS' + Math.random().toString(36).substr(2, 6).toUpperCase();
}

// Functions to handle number pad input for the amount
function addDigit(digit) {
    if (amountString.value.length >= 10) return; // Limit max length

    amountString.value += digit;
    updateAmountValue();
}

function addDoubleZero() {
    if (amountString.value.length + 2 > 10) return; // Check max length
    amountString.value += '00';
    updateAmountValue();
}

function addDecimal() {
    if (!amountString.value.includes('.')) {
        amountString.value += '.';
    }
}

function clearAmount() {
    amountString.value = '';
    amount.value = 0;
}

function deleteLastChar() {
    amountString.value = amountString.value.slice(0, -1);
    updateAmountValue();
}

// Update numeric amount based on amountString (input from number pad)
function updateAmountValue() {
    amount.value = parseFloat(amountString.value) || 0;
}

// Create an invoice when the Enter button is clicked
function enterAmount() {
    // Calculate the XMR equivalent using the exchange rate and convert to piconeros
    fetchExchangeRate()
        .then(() => {
            const xmrAmount = amount.value / exchangeRate.value; // Convert fiat to XMR if needed
            invoice.value.amount = convertToPiconeros(xmrAmount); // Convert XMR to piconeros
            console.log('amount=', amount.value, ' rate=', exchangeRate.value, ' xmr=', xmrAmount, ' pico', invoice.value.amount);

            invoice.value.order_id = generateOrderId();

            // Create the invoice
            return axios.post(
                `${import.meta.env.VITE_API_BASE}/merchant/invoice`,
                {
                    amount: invoice.value.amount,
                    order_id: invoice.value.order_id,
                    accept_url: invoice.value.accept_url,
                    cancel_url: invoice.value.cancel_url
                },
                { withCredentials: true }
            );
        })
        .then((response) => {
            // Redirect to the newly created invoice page
            window.location.href = `${import.meta.env.VITE_API_BASE}/p/page/${response.data.invoice_id}`;
        })
        .catch((error) => {
            if (error.response?.status === 401) {
                toast.add({ severity: 'error', summary: 'Not logged in', detail: error.code, life: 3000 });
                push({ path: '/auth/login' });
            } else {
                toast.add({
                    severity: 'error',
                    summary: 'Failed to create invoice',
                    detail: error.message,
                    life: 3000
                });
            }
        });
}
</script>

<template>
    <div>
        <div class="flex items-center justify-center">
            <div class="flex flex-col items-center justify-center">
                <div style="border-radius: 56px; padding: 0.3rem; background: linear-gradient(180deg, var(--primary-color) 10%, rgba(33, 150, 243, 0) 30%)">
                    <div class="w-full bg-surface-0 dark:bg-surface-900 py-20 px-8 sm:px-20" style="border-radius: 53px">
                        <div class="mb-5 flex gap-4">
                            <!-- InputText to display amount with selected currency symbol as suffix -->
                            <InputText v-model="amountString" size="large" :suffix="selectedCurrency" fluid readonly></InputText>
                            <Select size="large" v-model="selectedCurrency" :options="currencies"></Select>
                        </div>
                        <div class="grid grid-cols-3 gap-4 !text-2xl">
                            <Button size="large" label="1" @click="addDigit(1)"></Button>
                            <Button size="large" label="2" @click="addDigit(2)"></Button>
                            <Button size="large" label="3" @click="addDigit(3)"></Button>
                            <Button size="large" label="4" @click="addDigit(4)"></Button>
                            <Button size="large" label="5" @click="addDigit(5)"></Button>
                            <Button size="large" label="6" @click="addDigit(6)"></Button>
                            <Button size="large" label="7" @click="addDigit(7)"></Button>
                            <Button size="large" label="8" @click="addDigit(8)"></Button>
                            <Button size="large" label="9" @click="addDigit(9)"></Button>
                            <Button size="large" label="00" @click="addDoubleZero"></Button>
                            <Button size="large" label="0" @click="addDigit(0)"></Button>
                            <Button size="large" label="." @click="addDecimal"></Button>
                            <Button size="large" severity="danger" icon="pi pi-times-circle text-2xl" label=" " @click="clearAmount"></Button>
                            <Button size="large" severity="info" icon="pi pi-delete-left text-2xl" label=" " @click="deleteLastChar"></Button>
                            <Button size="large" severity="success" icon="pi pi-arrow-right text-2xl" label=" " @click="enterAmount"></Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
