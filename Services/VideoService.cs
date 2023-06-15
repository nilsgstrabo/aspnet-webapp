using Azure.Security.KeyVault.Secrets;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using System.Collections.Generic;
using System.Threading.Tasks;
using System.Linq;
namespace aspnet_webapp.Services {

    public class VideoInfo {
        public string Id { get; set; }
        public string FileName { get; set; }
    }

    public interface IVideoService {
        public IEnumerable<VideoInfo> GetVideos();
    }

    public class VideoService : IVideoService {
        private readonly IConfiguration _config;
        public VideoService(IConfiguration config)
        {
            _config=config;
        }

        public IEnumerable<VideoInfo> GetVideos()
        {
            var videos = new List<VideoInfo>();
            var rootPath=_config["VIDEO_PATH"];
            if(rootPath?.Length==0) {
                return videos;
            }
           
            var q=new Stack<string>(new string[]{rootPath});
            
            
            while (q.Count>0)
            {
                var dir=q.Pop();
                System.IO.Directory.GetDirectories(dir).ToList().ForEach(d=>q.Push(d));
                System.IO.Directory.GetFiles(dir).Select(f=>new System.IO.FileInfo(f)).ToList().ForEach(f=>videos.Add(new VideoInfo{FileName=f.FullName, Id=f.FullName.Replace('/','-')}));
            }
            return videos;

        }
    }

}