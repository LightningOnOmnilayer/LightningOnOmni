# Step by step to write your first lightning wallet



First, you should put all of [JS SDK](https://github.com/omnilaboratory/DebuggingTool/tree/master/sdk) files 
to your path of project.

And, following these steps below in your code:

* [Step 1: connect to an OBD node](#step-1-connect-to-an-obd-node)
* [Step 2: signup a new user](#step-2-signup-a-new-user)
* [Step 3: login using mnemonic words](#step-3-login-using-mnemonic-words)
* [Step 4: connect another user](#step-4-connect-another-user)
* [Step 5: open channel](#step-5-open-channel)
* [Step 6: create an invoice](#step-6-create-an-invoice)
* [Step 7: pay an invoice](#step-7-pay-an-invoice)
  

## Step 1: connect to an OBD node

<p align="center">
  <img width="750" alt="login" src="prototype/connectOBD.png">
</p>

Invoke **connectToServer** function from [wallet.js](https://github.com/omnilaboratory/DebuggingTool/blob/master/sdk/wallet.js) of SDK.

First parameter is `nodeAddress`. It's a URL of the OBD node like `ws://62.234.216.108:60020/wstest`.

Second parameter is `callback`. It's a callback function could be used to process the return data.

Third parameter is `globalCallback`. It's a callback function could be used to process the global messages.

#### Example Code:

```js
// URL of an OBD node
let nodeAddress = 'ws://62.234.216.108:60020/wstest';

// SDK API
connectToServer(nodeAddress, function(response) {
    console.info('Print the callback = ' + response);

    // Your code to process the callback data.
    // Example: Display the success or fail message on app screen.

}, function(globalResponse) {
    // Your code to process the global callback data.

});
```

Full example in GUI-tool you could be see [sdkConnect2OBD](https://github.com/omnilaboratory/DebuggingTool/blob/master/js/common.js) function.


## Step 2: signup a new user

<p align="center">
  <img width="750" alt="login" src="prototype/createAccount.png">
</p>

Invoke **genMnemonic** function from [wallet.js](https://github.com/omnilaboratory/DebuggingTool/blob/master/sdk/wallet.js) of SDK.

`genMnemonic` function is used to sign up a new user by hirarchecal deterministic wallet system integrated in the local client. Client generates mnemonic words and the hash of the mnemonic words as the UserID.

#### Example Code:

```js
// SDK API
let mnemonicWords = genMnemonic();

// Your code to process the data.
// Example: Display the mnemonic words on app screen.
```

Full example in GUI-tool you could be see [sdkGenMnemonic](https://github.com/omnilaboratory/DebuggingTool/blob/master/js/common.js) function.


## Step 3: login using mnemonic words

<p align="center">
  <img width="300" alt="login" src="prototype/loginWithSeed.png">
</p>

Invoke **logIn** function from [wallet.js](https://github.com/omnilaboratory/DebuggingTool/blob/master/sdk/wallet.js) of SDK.

First parameter is `mnemonic`. It's mnemonic words generated from **genMnemonic** function.

Second parameter is `callback`. It's a callback function could be used to process the return data.

#### Example Code:

```js
let mnemonic = 'mnemonic words';

// SDK API
logIn(mnemonic, function(e) {
    // Print the callback data
    console.info('SDK: -102001 logIn = ' + JSON.stringify(e));

    // SDK API: For auto pilot mode, register event for listening.
    // @param netType --> true: testnet  false: mainnet
    registerEvent(true);

    // Your code to process the callback data.
    // Example: Save mnemonic words used by a user to log in
    // SDK API
    saveMnemonic(e.userPeerId, mnemonic);
});
```

Full example in GUI-tool you could be see [sdkLogIn](https://github.com/omnilaboratory/DebuggingTool/blob/master/js/common.js) function.


## Step 4: connect another user

<p align="center">
  <img width="300" alt="login" src="prototype/connectRemoteUser.png">
</p>

Connect a remote user, then you will be able to creat your first channel in the network. If you already have a channel, you can jump to the next step to pay or create invoices to collect money.

Invoke **connectPeer** function from [wallet.js](https://github.com/omnilaboratory/DebuggingTool/blob/master/sdk/wallet.js) of SDK.

First parameter is `P2PPeer object`. It contains node address you will connect to. The node address returned by **logIn** function. It looks like this:

```
/ip4/62.234.216.108/tcp/3001/p2p/QmP1mQMzDRV2bKWdhwvPWFubSAz1gqJY44RjdYm3G5DFeF
```

Second parameter is `callback`. It's a callback function could be used to process the return data.

#### Example Code:

```js
let info                 = new P2PPeer();
info.remote_node_address = 'node address';

// SDK API
connectPeer(info, function(e) {
    // Print the callback data
    console.info('SDK: -102003 connectPeer = ' + JSON.stringify(e));

    // Your code to process the callback data.
});
```

Full example in GUI-tool you could be see [sdkConnectP2PPeer](https://github.com/omnilaboratory/DebuggingTool/blob/master/js/common.js) function.


## Step 5: open channel

Invoke **openChannel** function from [basic.js](https://github.com/omnilaboratory/DebuggingTool/blob/master/sdk/basic.js) of SDK.

First parameter is `myUserID`. It's the user id of the client currently logged in.

Second and third parameter are `nodeID` and `remoteUserID`. Both returned by **logIn** function. Explanation is [here](https://omnilaboratory.github.io/obd/#/GUI-tool?id=step-3-login-using-mnemonic-words)

Final parameter is `OpenChannelInfo object`. It contains `public key` of the address that you will use to create channel address. And `is_private` means the channel is public or private.

#### Example Code:

```js
let nodeID       = 'nodePeerId';
let remoteUserID = 'userPeerId';

let info            = new OpenChannelInfo();
info.funding_pubkey = 'public key of an address';
info.is_private     = true or false;

// SDK API
openChannel(myUserID, nodeID, remoteUserID, info);
```

Let's check out the `openChannel` function.

```js
function openChannel(myUserID, nodeID, remoteUserID, info) {
    obdApi.openChannel(nodeID, remoteUserID, info, function(e) {
        // Print the callback data
        console.info('SDK: -100032 openChannel = ' + JSON.stringify(e));

        // To simplify development, we save some data to local storage at client.
        // All of this is SDK APIs.
        saveCounterparties(myUserID, nodeID, remoteUserID);
        saveChannelID(e.temporary_channel_id);
        let privkey = getFundingPrivKeyFromPubKey(myUserID, info.funding_pubkey);
        saveFundingPrivKey(myUserID, e.temporary_channel_id, privkey, kTbFundingPrivKey);
    });
}
```

In `openChannel` function we call the OBD's function `openChannel`. The final parameter is `callback`. To simplify development, we save some data of callback and others to local storage at client.

Full example in GUI-tool you could be see [sdkOpenChannel](https://github.com/omnilaboratory/DebuggingTool/blob/master/js/common.js) function.


## Step 6: create an invoice 

<p align="center">
  <img width="300" alt="login" src="prototype/createInvoice.png">
</p>

Invoke **addInvoice** function from [htlc.js](https://github.com/omnilaboratory/DebuggingTool/blob/master/sdk/htlc.js) of SDK.

First parameter is `InvoiceInfo object`. Explanation is [here](https://omnilaboratory.github.io/obd/#/GUI-tool?id=step-6-create-an-invoice)

Second parameter is `callback`. It's a callback function could be used to process the return data.

#### Example Code:

```js
let info         = new InvoiceInfo();
info.property_id = 'property id of an asset';
info.amount      = 'amount';
info.h           = 'hash of r';
info.expiry_time = 'expiry time';
info.description = 'a memo';

// SDK API
addInvoice(info, function(e) {
    // Print the callback data
    console.info('SDK: -100402 addInvoice = ' + JSON.stringify(e));

    // Your code to process the callback data.
    // Example: create a QR code 
    // makeQRCode(e);
});
```

Full example in GUI-tool you could be see [sdkAddInvoice](https://github.com/omnilaboratory/DebuggingTool/blob/master/js/common.js) function.


## Step 7: pay an invoice 

<p align="center">
  <img width="750" alt="login" src="prototype/payInvoice.png">
</p>

Invoke **payInvoice** function from [htlc.js](https://github.com/omnilaboratory/DebuggingTool/blob/master/sdk/htlc.js) of SDK.

First parameter is `PayInvoiceInfo object`. 

Payment infomation is encoded in an invoice represented by a QR-code or a string. Payer either:
1）call **payInvoice** to pay the invoice created by payee,  
OR  
2）call **payInvoice** by filling all the payment items manuely. The example lists all the payment items.  

Second parameter is `callback`. It's a callback function could be used to process the return data.

#### Example Code:

```js
let info     = new PayInvoiceInfo();
let isInvPay = true or false;

if (isInvPay === true) { // pay an invoice
    info.invoice = 'invoice';
} else { // pay to a user id without an invoice
    info.recipient_user_peer_id = 'recipient user id';
    info.property_id = 'property id of an asset';
    info.amount      = 'amount';
    info.h           = 'hash of r';
    info.expiry_time = 'expiry time';
    info.description = 'a memo';
    info.is_private  = true or false;  // the channel is public or private
}

// SDK API
payInvoice(info);
```

Full example in GUI-tool you could be see [sdkPayInvoice](https://github.com/omnilaboratory/DebuggingTool/blob/master/js/common.js) function.
