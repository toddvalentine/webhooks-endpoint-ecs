<?php

$data_string = <<<BODY
[{
    "id": "9593",
    "webhook_trigger_type": "new_transaction",
    "user_id": "638",
    "cause_id": "4809",
    "transaction_type": "D,M",
    "transaction_id": "T1174265",
    "total_transaction_amount": "25.00",
    "transaction_tax_deductible_amount": "25.00",
    "fee_by_supporter": "Y",
    "processing_fee": "1.53",
    "total_charged_amount": "26.53",
    "payment_method": "Credit Card",
    "posted_amount": "25.00",
    "business_organization_name": "",
    "first_name": "Todd",
    "last_name": "Valentine",
    "net_received": "25.00",
    "campaign_name": "Multi Tool Campaign",
    "email": "donationsignaturevalidation@aol.com",
    "city": "Nashville",
    "state": "TN",
    "address": "123 Street ",
    "postal_code": "12345",
    "country": "US",
    "phone_no_type": "M",
    "country_code": "+1",
    "phone_number": "1234567890",
    "want_to_anonymous": "No",
    "individual_or_business": "Individual",
    "support_message": "",
    "contact_id": "CID638-1000056",
    "transaction_date": "2021-10-04 15:56:10",
    "settlement_date": "2021-10-04 15:56:10",
    "associated_activity_IDs": "D1129283_A,M1026828_A",
    "custom_field_status": "false"
}]
BODY;

var_dump($data_string);

$url = "Http://localhost:8080";

function create_hmac_sha256($requestbodydata){ 
    return base64_encode(hash_hmac('sha256',$requestbodydata, "whsec_jUWeYgMLFhIyw7EU", true));    
}

$get_hmac_sha256 = create_hmac_sha256($data_string);
$request_headers = array(
    "Content-Type: application/json",
    "X-vtypeio-Hmac-SHA256:" . $get_hmac_sha256                        
);
/*Y:: PD-1706 sending Hmac-SHA256 in headers end*/

// send parameter in post via curl call
$ch = curl_init($url);
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_POSTFIELDS, $data_string);
curl_setopt($ch, CURLOPT_HEADER, true);
//curl_setopt($ch, CURLOPT_HTTPHEADER, array('Content-Type: application/json'));
curl_setopt($ch, CURLOPT_HTTPHEADER, $request_headers); /*Y:: PD-1706 sending Hmac-SHA256*/
$response = curl_exec($ch);
curl_close($ch);

