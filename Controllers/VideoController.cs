using System.IO;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using Microsoft.AspNetCore.StaticFiles;
using Microsoft.Graph;
using System.Net.Http;
using System.Net.Http.Headers;
using Microsoft.Extensions.Hosting;
using Microsoft.AspNetCore.Http;
using aspnet_webapp.Services;
using System.Linq;
using Microsoft.AspNetCore.Routing;
using System.ComponentModel.DataAnnotations;
using System.Reflection.Metadata.Ecma335;
using System;
using Microsoft.Extensions.Caching.Memory;
using System.Text;

namespace aspnet_webapp.Controllers
{

    public class MegaStream : Stream
    {
        private readonly long size;
        private long pos;
        private readonly Stream data;
        private readonly ILogger _logger;

        public MegaStream(long size, ILogger logger)
        {
            var s=System.Linq.Enumerable.Repeat("this_is_the_content_of_the_megafile", 1000).Aggregate((a, s)=>a+s);
            data=new MemoryStream(Encoding.ASCII.GetBytes(s));
            this.size=size;
            _logger=logger;
        }

        public override bool CanRead => true;

        public override bool CanSeek => false;

        public override bool CanWrite => false;

        public override long Length => size;

        public override long Position { get => pos; set => pos=value; }

        public override void Flush()
        {
        }

        public override int Read(byte[] buffer, int offset, int count)
        {
            _logger?.LogInformation("Read data count: "+count.ToString());
            if (pos>=size) {
                return 0;
            }

            if (data.Position>=data.Length) {
                data.Position=0;
            }

            count = pos+count>size ? Convert.ToInt32(size-pos) : count;
            var bytesRead=data.Read(buffer,offset,count);
            
            pos+=bytesRead;

            return bytesRead;
            // if (pos>=size) {
            //     return 0;
            // }

            // int cnt = count>65_000 ? 65_000 : count;
            // int bytesRead = 0;
            // for (int i = 0; i < cnt; i++)
            // {
            //     if (pos+i>=size) {
            //         break;
            //     }
            //     buffer[i+offset]=Convert.ToByte((pos+i) % 26 +65);
            //     bytesRead++;
            // }

            // this.pos+=bytesRead;
            // return bytesRead;
        }

        public override long Seek(long offset, SeekOrigin origin)
        {
            throw new System.NotImplementedException();
        }

        public override void SetLength(long value)
        {
            throw new System.NotImplementedException();
        }

        public override void Write(byte[] buffer, int offset, int count)
        {
            throw new System.NotImplementedException();
        }
    }

    [ApiController]
    [Route("api/[controller]")]
    public class VideoController : ControllerBase
    {
        private readonly ILogger _logger;
        private readonly IConfiguration _configuration;
        private readonly IVideoService _videoService;

        public VideoController(IVideoService videoService, IConfiguration configuration,ILogger<VideoController> logger)
        {
            _logger = logger;
            _configuration = configuration;
            _videoService=videoService;
        }

        [HttpGet()]
        public IActionResult GetVideos()
        {
            var s=new FileStream("",FileMode.OpenOrCreate);
            return this.Ok();
        }


        [HttpGet("megafile")]
        public IActionResult GetMegaFile(string name)
        {
            long filesize=1_000_000;
            try
            {
                var s=Convert.ToInt64(Environment.GetEnvironmentVariable("MEGA_FILE_SIZE"));
                filesize = s >0 ? s : filesize;
            }
            catch (System.Exception ex)
            {
                _logger.LogError(ex, "failed to get mega file size, using default 1GB");
            }
            _logger.LogInformation("Stream megafile");
            return this.File(new MegaStream(filesize, _logger),"application/octet-stream", "megafile.txt");
           
        }

        [HttpGet("{name}")]
        public IActionResult GetVideo(string name)
        {
            _logger.LogInformation("Stream video {0}, range {1}", name, this.Request.GetTypedHeaders().Range?.ToString());
            var video=_videoService.GetVideos().FirstOrDefault(v=>v.Id==name);

            if(video==null) {
                return this.NotFound();
            }
            
            return this.PhysicalFile(video.FileName, "video/mp4", true);
        }

    }
}