export namespace manifest {
	
	export class Manifest {
	    name: string;
	    uid: string;
	    pubsub: string;
	    chiper: string;
	    optional: string;
	
	    static createFrom(source: any = {}) {
	        return new Manifest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.uid = source["uid"];
	        this.pubsub = source["pubsub"];
	        this.chiper = source["chiper"];
	        this.optional = source["optional"];
	    }
	}

}

export namespace models {
	
	export class Message {
	    datatype: number;
	    sender: string;
	    senderid: string;
	    data: string;
	    time: string;
	    sign: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.datatype = source["datatype"];
	        this.sender = source["sender"];
	        this.senderid = source["senderid"];
	        this.data = source["data"];
	        this.time = source["time"];
	        this.sign = source["sign"];
	    }
	}
	export class ProfileStorePublic {
	    type: string;
	    name: string;
	    id: string;
	    pub: string;
	    avatar: string;
	
	    static createFrom(source: any = {}) {
	        return new ProfileStorePublic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.name = source["name"];
	        this.id = source["id"];
	        this.pub = source["pub"];
	        this.avatar = source["avatar"];
	    }
	}

}

