<script setup>
import { stringifyExpression } from '@vue/compiler-core';

defineProps({
    mailTextbox: {
        From: String,
        To: String,
        Subject: String,
        Date: String,
        BodyMessage: String
    }
})
</script>

<template>
    <div class="sm:w-full m-4 relative overflow-x-auto shadow-md sm:rounded-lg border overflow-y-scroll">
        <p v-if = "mailTextbox!=null" class="mb-3 font-light text-gray-500 dark:text-gray-400 m-2 h-fill">
            <strong> FROM: </strong> {{ mailTextbox.From }} <br>
            <strong> TO: </strong> {{ mailTextbox.To }} <br>
            <strong> DATE: </strong> {{ mailTextbox.Date }} <br>
            <strong> SUBJECT: </strong> {{ mailTextbox.Subject }} <br> <br>
            {{ mailTextbox.BodyMessage }}
        </p>
    </div>
</template>